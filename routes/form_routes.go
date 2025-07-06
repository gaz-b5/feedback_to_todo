package routes

import (
	"David/formfunctions"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/apis"
)

// RegisterFormRoutes sets up custom API endpoints for form-related actions.
func RegisterFormRoutes(e *core.ServeEvent) {
	// Example: Register a POST route for form submission
	e.Router.POST("/api/forms/project/new", func(e *core.RequestEvent) error {
		return formfunctions.CreateProject(e)
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
