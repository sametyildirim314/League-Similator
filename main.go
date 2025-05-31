package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sametyildirim314/insider_case/config"
	"github.com/sametyildirim314/insider_case/database"
	"github.com/sametyildirim314/insider_case/routes"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()
	
	// Connect to database
	database.ConnectDB()
	
	// Initialize database schema
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Premier League Simulation",
	})
	
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())
	
	// Setup routes
	// NOTE: I kept it simple by moving all service code directly into controllers
	// This is not a best practice for larger applications, but it's easier for me to
	// understand as I'm learning Go. It avoids an extra layer of abstraction.
	routes.SetupTeamRoutes(app)
	routes.SetupMatchRoutes(app)
	routes.SetupLeagueRoutes(app)
	routes.SetupPredictionRoutes(app)
	routes.SetupSystemRoutes(app)
	
	// Add a simple health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Premier League Simulation API is running",
		})
	})
	
	// Start server
	port := ":" + strconv.Itoa(cfg.AppPort)
	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(port))
} 