package routes

import (
	"David/process_emails"
	"github.com/pocketbase/pocketbase/core"
)

//Email handler, supposed to be temporary till email recieving method is finalised
func RegisterEHandlerRoutes(e *core.ServeEvent) {
	
	e.Router.POST("/api/feedback/new", func(e *core.RequestEvent) error {
		return process_emails.ProcessEmailContent(e)
	})
	
}