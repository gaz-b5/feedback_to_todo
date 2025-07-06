package process_emails

import (
    "encoding/json"
    "net/http"
    "github.com/pocketbase/pocketbase/core"
)

type ProcessEmailInput struct {
    ProjectId string `json:"project_id"`
    Content   string `json:"content"`
}

func ProcessEmailContent(e *core.RequestEvent) error {
    // Parse input
    var input ProcessEmailInput
    if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
        return e.BadRequestError("Invalid request body", err)
    }
    if input.ProjectId == "" || input.Content == "" {
        return e.BadRequestError("Missing project_id or content", nil)
    }

    // Get the project record by ID
    project, err := e.App.FindRecordById("projects", input.ProjectId)
    if err != nil || project == nil {
        return e.BadRequestError("Project not found", err)
    }

    // Extract the project name
    projectName := project.GetString("name")

    // Find the emails collection
    collection, err := e.App.FindCollectionByNameOrId("emails")
    if err != nil {
        return e.InternalServerError("Emails collection not found", err)
    }

    // Create a new email record
    record := core.NewRecord(collection)
    record.Set("project", input.ProjectId) // assumes 'project' is a relation field in emails collection
    record.Set("content", input.Content)

    if err := e.App.Save(record); err != nil {
        return e.InternalServerError("Failed to save email content", err)
    }

    return e.JSON(http.StatusCreated, map[string]any{
        "message": "Email content processed and saved",
        "project_name": projectName,
        "email_record": record.PublicExport(),
    })
}