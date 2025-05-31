package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/controllers"
)

// SetupMatchRoutes sets up all routes for matches
func SetupMatchRoutes(app *fiber.App) {
	api := app.Group("/api")
	matches := api.Group("/matches")
	
	matches.Get("/", controllers.GetAllMatches)
	matches.Get("/week/:week", controllers.GetMatchesByWeek)
	matches.Post("/simulate/:week", controllers.SimulateWeek)
	matches.Post("/simulate-all", controllers.SimulateAllRemainingMatches)
} 