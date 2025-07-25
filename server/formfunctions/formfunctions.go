package formfunctions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

type ReomveUserByEmailInput struct {
	ProjectId string `json:"project_id"`
	Email     string `json:"email"`
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
	if input.ProjectId == "" || input.Email == "" || (input.Role != "admin" && input.Role != "viewer") {
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
	var input ReomveUserByEmailInput
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

	// 2. Find the user to remove by email
	userToRemove, err := e.App.FindAuthRecordByEmail("users", input.Email)
	if err != nil || userToRemove == nil {
		return e.BadRequestError("No user found with that email", err)
	}

	// 3. Check if the user is already a member
	filter = "project = {:project} && user_id = {:user_id}"
	params = map[string]any{
		"project": input.ProjectId,
		"user_id": userToRemove.Id,
	}
	membership, err := e.App.FindFirstRecordByFilter("users_projects", filter, params)
	if err != nil || membership == nil {
		return e.BadRequestError("User is not a member of this project", nil)
	}

	// 4. Remove the user from the project (delete the membership record)
	if err := e.App.Delete(membership); err != nil {
		return e.InternalServerError("Failed to remove user from project", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message":    "User removed from project successfully",
		"user_id":    userToRemove.Id,
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
	tasks, err := e.App.FindRecordsByFilter("tasks", "project = {:project}", "", 1000, 0, map[string]any{"project": projectId})
	if err != nil {
		return e.InternalServerError("Failed to fetch tasks", err)
	}

	// 4. Respond with the list of tasks
	return e.JSON(http.StatusOK, map[string]any{
		"tasks": tasks,
	})
}
