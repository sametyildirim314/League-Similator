package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/controllers"
)

// SetupLeagueRoutes sets up all routes for the league
func SetupLeagueRoutes(app *fiber.App) {
	api := app.Group("/api")
	league := api.Group("/league")
	
	league.Get("/table", controllers.GetLeagueTable)
	league.Get("/table/week/:week", controllers.GetLeagueTableForWeek)
} 