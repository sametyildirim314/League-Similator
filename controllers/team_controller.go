package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/database"
	"github.com/sametyildirim314/insider_case/models"
)

// GetAllTeams gets all teams from the database and returns them
func GetAllTeams(c *fiber.Ctx) error {
	// Fetch teams directly from database
	rows, err := database.DB.Query("SELECT id, name FROM teams")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get teams: " + err.Error(),
		})
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan team: " + err.Error(),
			})
		}
		teams = append(teams, team)
	}
	
	return c.JSON(teams)
}

// GetTeamByID gets a team by ID and returns it
func GetTeamByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}
	
	// Fetch team directly from database
	var team models.Team
	err = database.DB.QueryRow("SELECT id, name FROM teams WHERE id = $1", id).Scan(&team.ID, &team.Name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}
	
	return c.JSON(team)
} 