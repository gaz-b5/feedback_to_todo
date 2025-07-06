package formfunctions

import (
    "encoding/json"
    "net/http"

    "github.com/pocketbase/pocketbase/core"
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

    // Use core.NewRecord (not models.NewRecord)
    record := core.NewRecord(collection)
    record.Set("name", input.Name)
    record.Set("owner_id", user.Id) // or whatever field links to the user

    if err := e.App.Save(record); err != nil {
        return e.InternalServerError("Failed to create project", err)
    }

    // 5. Respond with success
    return e.JSON(http.StatusCreated, map[string]any{
        "message": "Project created successfully",
        "project": record.PublicExport(),
    })
}
