package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/controllers"
)

// SetupPredictionRoutes sets up all routes for predictions
func SetupPredictionRoutes(app *fiber.App) {
	api := app.Group("/api")
	predictions := api.Group("/predictions")
	
	predictions.Get("/", controllers.GetPredictions)
} 