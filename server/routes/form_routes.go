package routes

import (
	"David/formfunctions"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterFormRoutes sets up custom API endpoints for form-related actions.
func RegisterFormRoutes(e *core.ServeEvent) {
	e.Router.POST("/api/forms/project/new", func(e *core.RequestEvent) error {
		return formfunctions.CreateProject(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.DELETE("/api/forms/project/delete", func(e *core.RequestEvent) error {
		return formfunctions.DeleteProject(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.POST("/api/forms/project/addUser", func(e *core.RequestEvent) error {
		return formfunctions.AddUserToProjectByEmail(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.DELETE("/api/forms/project/removeUser", func(e *core.RequestEvent) error {
		return formfunctions.RemoveUserFromProject(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.GET("/api/forms/projects", func(e *core.RequestEvent) error {
		return formfunctions.GetProjects(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.GET("/api/forms/project/tasks", func(e *core.RequestEvent) error {
		return formfunctions.GetTasks(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.POST("/api/forms/project/task/setpriority", func(e *core.RequestEvent) error {
		return formfunctions.SetPriority(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.POST("/api/forms/project/task/setstatus", func(e *core.RequestEvent) error {
		return formfunctions.SetStatus(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.DELETE("/api/forms/project/task/delete", func(e *core.RequestEvent) error {
		return formfunctions.DeleteTask(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.POST("/api/forms/project/task/setprioritybulk", func(e *core.RequestEvent) error {
		return formfunctions.SetPriorityBulk(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.POST("/api/forms/project/task/setstatusbulk", func(e *core.RequestEvent) error {
		return formfunctions.SetStatusBulk(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.DELETE("/api/forms/project/task/deletebulk", func(e *core.RequestEvent) error {
		return formfunctions.DeleteTasksBulk(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.PATCH("/api/forms/project/task/update", func(e *core.RequestEvent) error {
		return formfunctions.UpdateTask(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	e.Router.PATCH("/api/forms/project/task/bulk", func(e *core.RequestEvent) error {
		return formfunctions.UpdateTasksBulk(e)
	}).Bind(apis.RequireAuth()) // Require authentication, remove if not needed

	// Example: Register a GET route for fetching form data
	// e.Router.GET("/api/forms/:id", func(e *core.RequestEvent) error {
	// 	return formfunctions.HandleGetForm(e)
	// }).Bind(apis.RequireAuth())

	// // Example: Register a DELETE route for deleting a form
	// e.Router.DELETE("/api/forms/:id", func(e *core.RequestEvent) error {
	// 	return formfunctions.HandleDeleteForm(e)
	// }).Bind(apis.RequireAuth())

	// Add more routes as needed...
}
