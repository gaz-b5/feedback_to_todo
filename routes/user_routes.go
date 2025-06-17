package routes

import (
	"David/controllers"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterUserRoutes(e *core.ServeEvent){
  users:=e.Router.Group("/user")
	users.POST("/register",controllers.RegisterUser)
	users.GET("/register",controllers.PathCheck)

	// users.GET("{id}",controllers.RegisterUserfunc)

}