package api

import (
	"github.com/sujayramesh/SMS/api/handlers"

	"fmt"

	"github.com/labstack/echo/v4"
)

// Function to declare all routes
func createRoutes(e *echo.Echo) {
	fmt.Println("Enter createRoutes")

	e.POST("/sessionManagement", handlers.CreateSession)
	e.DELETE("/sessionManagement", handlers.DeleteSession)
	e.PUT("/sessionManagement", handlers.ModifySession)
	e.GET("/sessionManagement", handlers.GetSession)

	// Init Session map the very first time.
	handlers.SessionList.Init()

	fmt.Println("Exit createRoutes")
}

//Function to create new echo object
func NewEcho() *echo.Echo {

	fmt.Println("Enter NewEcho")
	e := echo.New()

	//Set routes and handlers
	createRoutes(e)

	fmt.Println("Exit NewEcho")

	return e
}

func TriggerStaleSessionCleanup() {
	handlers.DeleteExpiredSessions()
}

func TriggerExitActions() {
	handlers.ResetSessions()
}
