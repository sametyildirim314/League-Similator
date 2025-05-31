package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/database"
)

// ResetSystem handles the request to reset the system
// This will clear all matches, predictions, and reset the league table
func ResetSystem(c *fiber.Ctx) error {
	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction: " + err.Error(),
		})
	}
	
	// Delete all matches
	_, err = tx.Exec("DELETE FROM matches")
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete matches: " + err.Error(),
		})
	}
	
	// Delete all predictions
	_, err = tx.Exec("DELETE FROM predictions")
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete predictions: " + err.Error(),
		})
	}
	
	// Reset league table
	_, err = tx.Exec(`
		UPDATE league_table SET 
		points = 0,
		played = 0,
		wins = 0,
		draws = 0,
		losses = 0,
		goals_for = 0,
		goals_against = 0,
		goal_difference = 0
	`)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reset league table: " + err.Error(),
		})
	}
	
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction: " + err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "System reset successful. All matches, predictions, and league table data have been cleared.",
	})
} 