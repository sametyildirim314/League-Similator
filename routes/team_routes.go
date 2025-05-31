package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/controllers"
)

// SetupTeamRoutes sets up all routes for teams
func SetupTeamRoutes(app *fiber.App) {
	api := app.Group("/api")
	teams := api.Group("/teams")
	
	teams.Get("/", controllers.GetAllTeams)
	teams.Get("/:id", controllers.GetTeamByID)
} 