package formfunctions

import (
    "encoding/json"
    "net/http"

    "github.com/pocketbase/pocketbase/core"
	"David/qdrant_api"
)

type ProjectInput struct {
    Name string `json:"name"`
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

type AddUserByEmailInput struct {
    ProjectId string `json:"project_id"`
    Email     string `json:"email"`
    Role      string `json:"role"` // "admin" or "viewer"
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
        "user_id":    requestingUser.Id,
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
    filter = "project = {:project} && user_id = {:user}"
    params = map[string]any{
        "project": input.ProjectId,
        "user_id":    userToAdd.Id,
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
        "message": "User added to project",
        "membership": joinRecord.PublicExport(),
    })
}
