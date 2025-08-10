package formfunctions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"David/data_point"
	"David/llm_functions"
	"David/qdrant_api"

	"github.com/pocketbase/pocketbase/core"
)

type ProjectInput struct {
	Name string `json:"name"`
}

type ProjectIdInput struct {
	Id string `json:"project_id"`
}

func CreateProject(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	var input ProjectInput
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.Name == "" {
		return e.BadRequestError("Project name is required", nil)
	}

	// 3. Check for uniqueness of project name for this user
	filter := "name = {:name}"
	params := map[string]any{"name": input.Name}
	existing, err := e.App.FindFirstRecordByFilter("projects", filter, params)
	if err == nil && existing != nil {
		return e.BadRequestError("Project name already exists", nil)
	}
	// You can check for not found error if needed

	// 4. Create and save the new project
	collection, err := e.App.FindCollectionByNameOrId("projects")
	if err != nil {
		return e.InternalServerError("Projects collection not found", err)
	}
	join_collection, err := e.App.FindCollectionByNameOrId("users_projects")
	if err != nil {
		return e.InternalServerError("users_projects collection not found", err)
	}

	// Use core.NewRecord (not models.NewRecord)
	record := core.NewRecord(collection)
	record.Set("name", input.Name)
	record.Set("owner_id", user.Id) // or whatever field links to the user

	if err := e.App.Save(record); err != nil {
		return e.InternalServerError("Failed to create project : projects", err)
	}

	join_record := core.NewRecord(join_collection)
	join_record.Set("user_id", user.Id) // or whatever field links to the user
	join_record.Set("project", record.Id)
	join_record.Set("role", "admin")

	if err := e.App.Save(join_record); err != nil {
		return e.InternalServerError("Failed to create project : users_projects", err)
	}

	//create qdrant collection
	qdrant_api.CreateCollection(input.Name)

	// 5. Respond with success
	return e.JSON(http.StatusCreated, map[string]any{
		"message": "Project created successfully",
		"project": record.PublicExport(),
	})
}

func DeleteProject(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	var input ProjectIdInput
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.Id == "" {
		return e.BadRequestError("Project Id is required", nil)
	}

	// 3. Find the project by Id
	project, err := e.App.FindRecordById("projects", input.Id)
	if project == nil {
		return e.BadRequestError("Project does not exist", nil)
	}

	// 4. Check if the user is the owner of the project
	if project.Get("owner_id") != user.Id {
		return e.ForbiddenError("Only the owner can delete the project", nil)
	}

	// 5. Delete all user-project memberships for this project
	memberships, err := e.App.FindRecordsByFilter(
		"users_projects",
		"project = {:project}",
		"",
		1,
		1000,
		map[string]any{"project": project.Id},
	)
	if err == nil && memberships != nil {
		for _, m := range memberships {
			_ = e.App.Delete(m)
		}
	}

	// 6. Delete the project record
	if err := e.App.Delete(project); err != nil {
		return e.InternalServerError("Failed to delete project", err)
	}

	// 7. Optionally: delete the Qdrant collection
	qdrant_api.DeleteCollection(project.Get("name").(string))

	// 8. Respond with success
	return e.JSON(http.StatusOK, map[string]any{
		"message": "Project deleted successfully",
		"project": project.PublicExport(),
	})
}

type AddUserByEmailInput struct {
	ProjectId string `json:"project_id"`
	Email     string `json:"email"`
	Role      string `json:"role"` // "admin" or "viewer"
}

type ReomveUserByIdInput struct {
	ProjectId string `json:"project"`
	MemberId  string `json:"member_id"`
}

// Handler: Only allow if the requester is an admin for the project
func AddUserToProjectByEmail(e *core.RequestEvent) error {
	requestingUser := e.Auth
	if requestingUser == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// Parse input
	var input AddUserByEmailInput
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.ProjectId == "" || input.Email == "" {
		return e.BadRequestError("Missing or invalid fields", nil)
	}

	// 1. Check if the requesting user is an admin for this project
	filter := "project = {:project} && user_id = {:user_id} && role = {:role}"
	params := map[string]any{
		"project": input.ProjectId,
		"user_id": requestingUser.Id,
		"role":    "admin",
	}
	adminRecord, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if err != nil || adminRecord == nil {
		return e.ForbiddenError("You are not an admin for this project", nil)
	}

	// 2. Find the user to add by email
	userToAdd, err := e.App.FindAuthRecordByEmail("users", input.Email)
	if err != nil || userToAdd == nil {
		return e.BadRequestError("No user found with that email", err)
	}

	// 3. Check if the user is already a member
	filter = "project = {:project} && user_id = {:user_id}"
	params = map[string]any{
		"project": input.ProjectId,
		"user_id": userToAdd.Id,
	}
	existing, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if existing != nil {
		return e.BadRequestError("User is already a member of this project", nil)
	}

	// 4. Add the user to the project
	joinCollection, err := e.App.FindCollectionByNameOrId("users_projects")
	if err != nil {
		return e.InternalServerError("users_projects collection not found", err)
	}
	joinRecord := core.NewRecord(joinCollection)
	joinRecord.Set("project", input.ProjectId)
	joinRecord.Set("user_id", userToAdd.Id)
	joinRecord.Set("role", input.Role)

	if err := e.App.Save(joinRecord); err != nil {
		return e.InternalServerError("Failed to add user to project", err)
	}

	return e.JSON(http.StatusCreated, map[string]any{
		"message":    "User added to project",
		"membership": joinRecord.PublicExport(),
	})
}

func RemoveUserFromProject(e *core.RequestEvent) error {
	requestingUser := e.Auth
	if requestingUser == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// Parse input
	var input ReomveUserByIdInput
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.ProjectId == "" || input.MemberId == "" {
		return e.BadRequestError("Missing or invalid fields", nil)
	}

	// 1. Check if the requesting user is an admin for this project
	filter := "project = {:project} && user_id = {:user_id} && role = {:role}"
	params := map[string]any{
		"project": input.ProjectId,
		"user_id": requestingUser.Id,
		"role":    "admin",
	}
	adminRecord, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if err != nil || adminRecord == nil {
		return e.ForbiddenError("You are not an admin for this project", nil)
	}

	// 3. Check if the user is already a member
	filter = "project = {:project} && user_id = {:user_id}"
	params = map[string]any{
		"project": input.ProjectId,
		"user_id": input.MemberId,
	}
	membership, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	project, err := e.App.FindRecordById("projects", input.ProjectId)
	if err != nil || membership == nil {
		return e.BadRequestError("User is not a member of this project", nil)
	} else if membership.GetString("user_id") == requestingUser.Id {
		return e.BadRequestError("You cannot remove yourself from the project", nil)
	} else if membership.GetString("role") == project.GetString("owner_id") {
		return e.BadRequestError("You cannot remove the project owner from the project", nil)
	}

	// 4. Remove the user from the project (delete the membership record)
	if err := e.App.Delete(membership); err != nil {
		return e.InternalServerError("Failed to remove user from project", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message":    "User removed from project successfully",
		"user_id":    input.MemberId,
		"project_id": input.ProjectId,
	})

}

func GetProjects(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Find all projects where the user is a member
	filter := "user_id = {:user_id}"
	params := map[string]any{
		"user_id": user.Id,
	}
	memberships, err := e.App.FindRecordsByFilter("users_projects", filter, "", 1000, 0, params)
	if err != nil {
		return e.InternalServerError("Failed to fetch projects", err)
	}

	// 3. Collect project IDs from memberships
	projectIds := make([]string, len(memberships))
	for i, m := range memberships {
		projectIds[i] = m.GetString("project")
	}

	// 4. Fetch project details
	if len(projectIds) == 0 {
		return e.JSON(http.StatusOK, map[string]any{"projects": []any{}})
	}

	var filters []string
	params = make(map[string]any)
	for i, id := range projectIds {
		key := fmt.Sprintf("id%d", i)
		filters = append(filters, fmt.Sprintf("id = {:%s}", key)) // use equality check
		params[key] = id
	}

	filter = strings.Join(filters, " || ")
	projects, err := e.App.FindRecordsByFilter("projects", filter, "", 1000, 0, params)
	if err != nil {
		return e.InternalServerError("Failed to fetch projects", err)
	}

	// 5. Respond with the list of projects
	return e.JSON(http.StatusOK, map[string]any{
		"projects": projects,
	})
}

func GetTasks(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	projectId := e.Request.URL.Query().Get("projectId")
	if projectId == "" {
		return e.BadRequestError("Project Id is required", nil)
	}

	// get all the tasks for the project, if the user is a member of the project
	filter := "project = {:project} && user_id = {:user_id}"
	params := map[string]any{
		"project": projectId,
		"user_id": user.Id,
	}

	memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if memberRecord == nil {
		return e.ForbiddenError("You are not a member for this project", nil)
	}

	// 3. Find all tasks for the project
	tasks, err := e.App.FindRecordsByFilter("tasks", "project = {:project}", "-created", 1000, 0, map[string]any{"project": projectId})
	if err != nil {
		return e.InternalServerError("Failed to fetch tasks", err)
	}

	// 4. Respond with the list of tasks
	return e.JSON(http.StatusOK, map[string]any{
		"tasks": tasks,
	})
}

func DeleteTask(e *core.RequestEvent) error {
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	var input struct {
		TaskId string `json:"task_id"`
	}
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.TaskId == "" {
		return e.BadRequestError("Invalid task ID", nil)
	}

	// 3. Find the task by ID
	task, _ := e.App.FindRecordById("tasks", input.TaskId)
	if task == nil {
		return e.BadRequestError("Task does not exist", nil)
	}

	// 4. Check if the user is a member of the project
	filter := "project = {:project} && user_id = {:user_id}"
	params := map[string]any{
		"project": task.GetString("project"),
		"user_id": user.Id,
	}
	memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if memberRecord == nil {
		return e.ForbiddenError("You are not a member for this project", nil)
	}

	projectId := task.GetString("project")

	// 1. Fetch the project record by projectId
	project, err := e.App.FindRecordById("projects", projectId)
	if err != nil {
		return e.InternalServerError("Failed to find project", err)
	}
	if project == nil {
		return e.BadRequestError("Project does not exist for the task", nil)
	}

	// 2. Extract project name from the project record
	projectName := project.GetString("name")
	if projectName == "" {
		return e.InternalServerError("Project name is empty", nil)
	}

	// Then use projectName as Qdrant collection name
	collectionName := projectName

	qdrantDBId, ok := task.Get("qdrantDBId").(string)
	if !ok || qdrantDBId == "" {
		return e.InternalServerError("Task missing qdrantDBId for Qdrant", nil)
	}

	// Create point ID for Qdrant deletion (assuming UUID strings used)
	err = qdrant_api.DeleteTaskQdrant(collectionName, qdrantDBId)
	if err != nil {
		return e.InternalServerError("Failed to delete task from Qdrant", err)
	}

	// 5. Delete the task
	if err := e.App.Delete(task); err != nil {
		return e.InternalServerError("Failed to delete task", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "Task deleted successfully",
	})

}

func DeleteTasksBulk(e *core.RequestEvent) error {
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// Parse JSON input expecting an array of task IDs
	var input struct {
		TaskIds []string `json:"task_ids"`
	}
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if len(input.TaskIds) == 0 {
		return e.BadRequestError("No task IDs provided", nil)
	}

	// Find the first task to get the project and collection name
	firstTask, err := e.App.FindRecordById("tasks", input.TaskIds[0])
	if err != nil {
		return e.InternalServerError("Failed to find first task", err)
	}
	if firstTask == nil {
		return e.BadRequestError("First task does not exist", nil)
	}

	// Fetch project record once from first task's project
	projectId := firstTask.GetString("project")
	project, err := e.App.FindRecordById("projects", projectId)
	if err != nil {
		return e.InternalServerError("Failed to find project", err)
	}
	if project == nil {
		return e.BadRequestError("Project does not exist for the first task", nil)
	}

	collectionName := project.GetString("name")
	if collectionName == "" {
		return e.InternalServerError("Project name is empty", nil)
	}

	var qdrantIDs []string
	var tasksToDelete []*core.Record

	for _, taskId := range input.TaskIds {
		task, err := e.App.FindRecordById("tasks", taskId)
		if err != nil {
			return e.InternalServerError("Failed to find task", err)
		}
		if task == nil {
			return e.BadRequestError(fmt.Sprintf("Task %s does not exist", taskId), nil)
		}

		// Check project consistency
		if task.GetString("project") != projectId {
			return e.BadRequestError(fmt.Sprintf("Task %s does not belong to the same project", taskId), nil)
		}

		// Check membership
		filter := "project = {:project} && user_id = {:user_id}"
		params := map[string]any{
			"project": projectId,
			"user_id": user.Id,
		}
		memberRecord, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
		if err != nil {
			return e.InternalServerError("Failed to verify project membership", err)
		}
		if memberRecord == nil {
			return e.ForbiddenError(fmt.Sprintf("You are not a member of the project for task %s", taskId), nil)
		}

		qdrantDBId, ok := task.Get("qdrantDBId").(string)
		if !ok || qdrantDBId == "" {
			return e.InternalServerError(fmt.Sprintf("Task %s missing uniqueDBId for Qdrant", taskId), nil)
		}

		qdrantIDs = append(qdrantIDs, qdrantDBId)
		tasksToDelete = append(tasksToDelete, task)
	}

	// Bulk delete from Qdrant using single collection name
	err = qdrant_api.DeleteTasksQdrantBulk(collectionName, qdrantIDs)
	if err != nil {
		return e.InternalServerError("Failed to delete tasks from Qdrant", err)
	}

	// Delete all app tasks
	for _, task := range tasksToDelete {
		if err := e.App.Delete(task); err != nil {
			return e.InternalServerError("Failed to delete task", err)
		}
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": fmt.Sprintf("%d tasks deleted successfully", len(tasksToDelete)),
	})
}

func UpdateTask(e *core.RequestEvent) error {
	// 1. Get the authenticated user
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input with pointer fields for optional detection
	var input struct {
		TaskId   string  `json:"task_id"`
		Status   *string `json:"status,omitempty"`
		Priority *string `json:"priority,omitempty"`
		Nature   *string `json:"nature,omitempty"`
		Assigned *string `json:"assigned,omitempty"`
	}
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.TaskId == "" {
		return e.BadRequestError("Task ID is required", nil)
	}
	// At least one field must be present to update
	if input.Status == nil && input.Priority == nil && input.Nature == nil && input.Assigned == nil {
		return e.BadRequestError("No fields to update", nil)
	}

	// 3. Find the task by ID
	task, _ := e.App.FindRecordById("tasks", input.TaskId)
	if task == nil {
		return e.BadRequestError("Task does not exist", nil)
	}

	// 4. Check if the user is a member of the project
	filter := "project = {:project} && user_id = {:user_id}"
	params := map[string]any{
		"project": task.GetString("project"),
		"user_id": user.Id,
	}
	memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if memberRecord == nil {
		return e.ForbiddenError("You are not a member of this project", nil)
	}

	// 5. Update only provided fields
	if input.Status != nil {
		task.Set("status", *input.Status)
	}
	if input.Priority != nil {
		task.Set("priority", *input.Priority)
	}
	if input.Nature != nil {
		task.Set("nature", *input.Nature)
	}
	if input.Assigned != nil {
		task.Set("assigned", *input.Assigned)
	}

	// 6. Save the task
	if err := e.App.Save(task); err != nil {
		return e.InternalServerError("Failed to update task", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "Task updated successfully",
	})
}

func UpdateTasksBulk(e *core.RequestEvent) error {
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// Accept fields as pointers to detect their presence (not just empty string)
	var input struct {
		Tasks []struct {
			TaskId string `json:"task_id"`
		} `json:"tasks"`
		Status   *string `json:"status,omitempty"`
		Priority *string `json:"priority,omitempty"`
		Nature   *string `json:"nature,omitempty"`
	}

	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	if len(input.Tasks) == 0 {
		return e.BadRequestError("No tasks provided", nil)
	}

	// Do not require at least one field, but optionally you may check that at least one is set
	if input.Status == nil && input.Priority == nil && input.Nature == nil {
		return e.BadRequestError("No fields to update", nil)
	}

	for _, taskInfo := range input.Tasks {
		if taskInfo.TaskId == "" {
			return e.BadRequestError("Invalid task ID", nil)
		}

		task, _ := e.App.FindRecordById("tasks", taskInfo.TaskId)
		if task == nil {
			return e.BadRequestError("Task does not exist: "+taskInfo.TaskId, nil)
		}

		filter := "project = {:project} && user_id = {:user_id}"
		params := map[string]any{
			"project": task.GetString("project"),
			"user_id": user.Id,
		}
		memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
		if memberRecord == nil {
			return e.ForbiddenError("You are not member of the project: "+task.GetString("project"), nil)
		}

		// Only set fields if present
		if input.Status != nil {
			task.Set("status", *input.Status)
		}
		if input.Priority != nil {
			task.Set("priority", *input.Priority)
		}
		if input.Nature != nil {
			task.Set("nature", *input.Nature)
		}

		if err := e.App.Save(task); err != nil {
			return e.InternalServerError("Failed to update task "+taskInfo.TaskId, err)
		}
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "Tasks updated successfully",
	})
}

func GetMemebers(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	projectId := e.Request.URL.Query().Get("projectId")
	if projectId == "" {
		return e.BadRequestError("Project Id is required", nil)
	}

	// 3. Find the project by Id
	project, err := e.App.FindRecordById("projects", projectId)
	if project == nil {
		return e.BadRequestError("Project does not exist", nil)
	}

	// 4. Check if the user is the member of the project
	filter := "project = {:project} && user_id = {:user_id}"
	params := map[string]any{
		"project": projectId,
		"user_id": user.Id,
	}
	memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if memberRecord == nil {
		return e.ForbiddenError("You are not member of the project: "+projectId, nil)
	}

	// 5. Find all users_projects records for this project
	filter = "project = {:project}"
	params = map[string]any{"project": projectId}
	memberships, err := e.App.FindRecordsByFilter("users_projects", filter, "", 1000, 0, params)
	if err != nil {
		return e.InternalServerError("Failed to fetch members", err)
	}

	// 6. Collect user IDs and roles from memberships
	members := make([]map[string]any, len(memberships))
	for i, m := range memberships {
		user, _ := e.App.FindRecordById("users", m.GetString("user_id"))

		members[i] = map[string]any{
			"user_id": m.GetString("user_id"),
			"email":   user.GetString("email"),
			"name":    user.GetString("name"),
			"role":    m.GetString("role"),
		}
	}

	// 7. Respond with the list of members
	return e.JSON(http.StatusOK, map[string]any{
		"members": members,
	})
}

// For when the user wants to add a task directly, without any llm intervention
func AddTaskDirect(e *core.RequestEvent) error {
	// 1. Get the authenticated user (now e.Auth)
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// 2. Parse the JSON input
	var input struct {
		ProjectId   string `json:"project_id"`
		Description string `json:"description"`
		Nature      string `json:"nature"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// 3. Check if the user is a member of the project

	//get the project by ID
	project, err := e.App.FindRecordById("projects", input.ProjectId)
	if project == nil || err != nil {
		return e.BadRequestError("Project does not exist", nil)
	}

	filter := "project = {:project} && user_id = {:user_id}"
	params := map[string]any{
		"project": input.ProjectId,
		"user_id": user.Id,
	}
	memberRecord, _ := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if memberRecord == nil {
		return e.ForbiddenError("You are not a member for this project", nil)
	}

	// 4. Create and save the new task
	taskCollection, err := e.App.FindCollectionByNameOrId("tasks")
	if err != nil {
		return e.InternalServerError("Tasks collection not found", err)
	}

	task := core.NewRecord(taskCollection)
	task.Set("title", llm_functions.GenerateTitle(input.Description))
	task.Set("description", input.Description)
	task.Set("occurrence", 1)
	task.Set("priority", input.Priority)
	task.Set("nature", input.Nature)
	task.Set("status", input.Status)
	task.Set("project", input.ProjectId)

	if err := e.App.Save(task); err != nil {
		return e.InternalServerError("Failed to create task", err)
	}

	dataPoint := data_point.DataPoint{
		TaskID:    task.Id,
		Embedding: llm_functions.GenerateEmbedding(input.Description),
	}

	projectName, _ := project.Get("name").(string)
	qdrantDBId := qdrant_api.UpdateAndCreateDataPoint(dataPoint, projectName)

	if qdrantDBId == "" {
		return e.InternalServerError("Failed to create task in Qdrant", nil)
	}

	task.Set("qdrantDBId", qdrantDBId)

	if err := e.App.Save(task); err != nil {
		return e.InternalServerError("Failed to update qdrantDBId in task", err)
	}

	joinCollection, err := e.App.FindCollectionByNameOrId("emails_tasks")
	if err != nil {
		return e.InternalServerError("emails_tasks collection not found", err)
	}
	joinRecord := core.NewRecord(joinCollection)
	joinRecord.Set("task", task.Id)

	if err := e.App.Save(joinRecord); err != nil {
		return e.InternalServerError("Failed to add email to task", err)
	}

	return e.JSON(http.StatusCreated, map[string]any{
		"message": "Task created successfully",
		"task":    task.PublicExport(),
	})
}

type UpdateUserRoleInput struct {
	ProjectId string `json:"project_id"`
	MemberId  string `json:"member_id"`
	Role      string `json:"role"` // "admin" or "viewer" (adjust to your allowed set)
}

func UpdateUserRoleInProject(e *core.RequestEvent) error {
	user := e.Auth
	if user == nil {
		return e.UnauthorizedError("Not authenticated", nil)
	}

	// Parse input
	var input UpdateUserRoleInput
	if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}
	if input.ProjectId == "" || input.MemberId == "" {
		return e.BadRequestError("Missing or invalid fields", nil)
	}

	// 1. Check if requester is an admin of the project
	filter := "project = {:project} && user_id = {:user_id} && role = {:role}"
	params := map[string]any{
		"project": input.ProjectId,
		"user_id": user.Id,
		"role":    "admin",
	}
	adminRecord, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if err != nil || adminRecord == nil {
		return e.ForbiddenError("You are not an admin for this project", nil)
	}

	// 2. Get the project and check if MemberId is the owner
	project, err := e.App.FindRecordById("projects", input.ProjectId)
	if err != nil || project == nil {
		return e.BadRequestError("Project does not exist", nil)
	}
	ownerId := project.GetString("owner_id")
	if ownerId == input.MemberId {
		return e.BadRequestError("Cannot change owner's role", nil)
	}

	// 3. Get the membership record to update
	filter = "project = {:project} && user_id = {:user_id}"
	params = map[string]any{
		"project": input.ProjectId,
		"user_id": input.MemberId,
	}
	memberRecord, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if err != nil || memberRecord == nil {
		return e.BadRequestError("User is not a member of this project", nil)
	}

	// 4. Update the role
	memberRecord.Set("role", input.Role)
	if err := e.App.Save(memberRecord); err != nil {
		return e.InternalServerError("Failed to update user role", err)
	}

	// 5. Response
	return e.JSON(http.StatusOK, map[string]any{
		"message":   "User role updated",
		"member_id": input.MemberId,
		"role":      input.Role,
	})
}
