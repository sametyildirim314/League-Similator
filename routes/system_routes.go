package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/controllers"
)

// SetupSystemRoutes sets up all system-related routes
func SetupSystemRoutes(app *fiber.App) {
	api := app.Group("/api")
	system := api.Group("/system")
	
	system.Post("/reset", controllers.ResetSystem)
} 