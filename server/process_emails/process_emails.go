package process_emails

import (
	"David/data_point"
	"David/llm_functions"
	"David/qdrant_api"
	"encoding/json"
	"net/http"
	"strings"
	"time"

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

	// Generate tasks as a single string in a specified format from the llm
	tasks := llm_functions.GetTasks(input.Content)

	for _, task := range tasks {
		isBug := false
		taskQuery := strings.Split(task, "=")

		if taskQuery[1] == "true" {
			isBug = true
		}

		dataPoint := data_point.DataPoint{
			Content:   taskQuery[0],
			IsBug:     isBug,
			RepCount:  1,
			Priority:  1,
			Timestamp: time.Now(),
			Status:    data_point.Pending,
		}

		dataPoint.Embedding = llm_functions.GenerateEmbedding(dataPoint.Content)

		limit := uint64(1)

		// Fetch the project record by ID
		project, err := e.App.FindRecordById("projects", input.ProjectId)
		if err != nil || project == nil {
			return e.BadRequestError("Project not found", err)
		}

		// Get the project name
		projectName := project.GetString("name")

		// Use projectName as the Qdrant collection name
		TaskID, numResults := qdrant_api.ReturnClosestTaskID(projectName, limit, dataPoint.Embedding)

		if numResults > 0 {
			taskRecord, err := e.App.FindRecordById("tasks", TaskID)
			if err != nil {
				return e.BadRequestError("Task Record not found", err)
			}
			description := taskRecord.GetString("description")

			if llm_functions.CompareStrings(taskQuery[0], description) {
				// Get the current occurrence value
				current := taskRecord.GetInt("occurrence") // Use GetFloat if your field is float

				// Increment by one
				taskRecord.Set("occurrence", current+1)

				// Update the record
				if err := e.App.Save(taskRecord); err != nil {
					e.BadRequestError("Could not update taskRecord", err)
				}

				joinCollection, err := e.App.FindCollectionByNameOrId("emails_tasks")
				if err != nil {
					return e.InternalServerError("emails_tasks collection not found", err)
				}
				joinRecord := core.NewRecord(joinCollection)
				joinRecord.Set("email", record.Id)
				joinRecord.Set("task", taskRecord.Id)

				if err := e.App.Save(joinRecord); err != nil {
					return e.InternalServerError("Failed to add email to task", err)
				}

			} else {
				tasksCollection, err := e.App.FindCollectionByNameOrId("tasks")
				if err != nil {
					return e.InternalServerError("Tasks collection not found", err)
				}

				newTask := core.NewRecord(tasksCollection)
				newTask.Set("description", dataPoint.Content)
				if dataPoint.IsBug {
					newTask.Set("nature", "bug")
				} else {
					newTask.Set("nature", "feature")
				}
				newTask.Set("occurrence", dataPoint.RepCount)
				newTask.Set("priority", dataPoint.Priority)
				// newTask.Set("timestamp", dataPoint.Timestamp)
				newTask.Set("status", dataPoint.Status.String())
				// newTask.Set("embedding", dataPoint.Embedding)
				newTask.Set("project", input.ProjectId)
				newTask.Set("title", llm_functions.GenerateTitle(dataPoint.Content))

				if err := e.App.Save(newTask); err != nil {
					return e.InternalServerError("Could not create new task", err)
				}
				dataPoint.TaskID = newTask.Id
				uniqueDBId := qdrant_api.UpdateAndCreateDataPoint(dataPoint, projectName)

				newTask.Set("uniqueDBId", uniqueDBId)

				// Save the task record again to update uniqueDBId field
				if err := e.App.Save(newTask); err != nil {
					return e.InternalServerError("Failed to update uniqueDBId in task", err)
				}

				joinCollection, err := e.App.FindCollectionByNameOrId("emails_tasks")
				if err != nil {
					return e.InternalServerError("emails_tasks collection not found", err)
				}
				joinRecord := core.NewRecord(joinCollection)
				joinRecord.Set("email", record.Id)
				joinRecord.Set("task", newTask.Id)

				if err := e.App.Save(joinRecord); err != nil {
					return e.InternalServerError("Failed to add email to task", err)
				}
			}
		} else {
			tasksCollection, err := e.App.FindCollectionByNameOrId("tasks")
			if err != nil {
				return e.InternalServerError("Tasks collection not found", err)
			}

			newTask := core.NewRecord(tasksCollection)
			newTask.Set("description", dataPoint.Content)
			if dataPoint.IsBug {
				newTask.Set("nature", "bug")
			} else {
				newTask.Set("nature", "feature")
			}
			newTask.Set("occurrence", dataPoint.RepCount)
			newTask.Set("priority", dataPoint.Priority)
			// newTask.Set("timestamp", dataPoint.Timestamp)
			newTask.Set("status", dataPoint.Status.String())
			// newTask.Set("embedding", dataPoint.Embedding)
			newTask.Set("project", input.ProjectId)
			newTask.Set("title", llm_functions.GenerateTitle(dataPoint.Content))

			if err := e.App.Save(newTask); err != nil {
				return e.InternalServerError("Could not create new task", err)
			}
			dataPoint.TaskID = newTask.Id
			uniqueDBId := qdrant_api.UpdateAndCreateDataPoint(dataPoint, projectName)

			newTask.Set("uniqueDBId", uniqueDBId)

			// Save the task record again to update uniqueDBId field
			if err := e.App.Save(newTask); err != nil {
				return e.InternalServerError("Failed to update uniqueDBId in task", err)
			}

			joinCollection, err := e.App.FindCollectionByNameOrId("emails_tasks")
			if err != nil {
				return e.InternalServerError("emails_tasks collection not found", err)
			}
			joinRecord := core.NewRecord(joinCollection)
			joinRecord.Set("email", record.Id)
			joinRecord.Set("task", newTask.Id)

			if err := e.App.Save(joinRecord); err != nil {
				return e.InternalServerError("Failed to add email to task", err)
			}

		}
	}

	return e.JSON(http.StatusCreated, map[string]any{
		"message":      "Email content processed and saved",
		"project_name": projectName,
		"email_record": record.PublicExport(),
	})
}
