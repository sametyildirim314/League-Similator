package controllers

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/database"
	"github.com/sametyildirim314/insider_case/models"
)

// GetPredictions handles the request to get all predictions
func GetPredictions(c *fiber.Ctx) error {
	// Query predictions directly from the database
	query := `
		SELECT p.id, p.team_id, p.predicted_position, p.predicted_points, 
		       p.prediction_percentage, p.created_at, t.id, t.name
		FROM predictions p
		JOIN teams t ON p.team_id = t.id
		ORDER BY p.predicted_position
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get predictions: " + err.Error(),
		})
	}
	defer rows.Close()
	
	var predictions []models.Prediction
	for rows.Next() {
		var prediction models.Prediction
		var createdAt sql.NullTime
		
		err := rows.Scan(
			&prediction.ID, &prediction.TeamID, &prediction.PredictedPosition,
			&prediction.PredictedPoints, &prediction.PredictionPercentage, &createdAt,
			&prediction.Team.ID, &prediction.Team.Name,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan prediction: " + err.Error(),
			})
		}
		
		if createdAt.Valid {
			prediction.CreatedAt = createdAt.Time
		}
		
		predictions = append(predictions, prediction)
	}
	
	return c.JSON(predictions)
}

// SubmitPrediction handles the request to submit a new prediction
func SubmitPrediction(c *fiber.Ctx) error {
	var prediction models.Prediction
	
	if err := c.BodyParser(&prediction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}
	
	// Insert prediction directly into database
	query := `
		INSERT INTO predictions (team_id, predicted_position, predicted_points, prediction_percentage)
		VALUES ($1, $2, $3, $4)
	`
	
	_, err := database.DB.Exec(
		query, 
		prediction.TeamID, 
		prediction.PredictedPosition, 
		prediction.PredictedPoints, 
		prediction.PredictionPercentage,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit prediction: " + err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Prediction submitted successfully",
		"prediction": prediction,
	})
}

// GenerateChampionshipProbabilities handles the request to generate prediction percentages
func GenerateChampionshipProbabilities(c *fiber.Ctx) error {
	// Get current league standings from database
	standingsQuery := `
		SELECT lt.points, lt.played, lt.wins, lt.draws, lt.losses, 
		       lt.goals_for, lt.goals_against, lt.goal_difference,
		       t.id, t.name
		FROM league_table lt
		JOIN teams t ON lt.team_id = t.id
		ORDER BY lt.points DESC, lt.goal_difference DESC, lt.goals_for DESC
	`
	
	rows, err := database.DB.Query(standingsQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get league standings: " + err.Error(),
		})
	}
	
	var standings []models.TeamStats
	for rows.Next() {
		var stats models.TeamStats
		err := rows.Scan(
			&stats.Points, &stats.Played, &stats.Wins, &stats.Draws, &stats.Losses,
			&stats.GoalsFor, &stats.GoalsAgainst, &stats.GoalDifference,
			&stats.Team.ID, &stats.Team.Name,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan team stats: " + err.Error(),
			})
		}
		
		standings = append(standings, stats)
	}
	rows.Close()
	
	// Clear existing predictions
	_, err = database.DB.Exec("DELETE FROM predictions")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear existing predictions: " + err.Error(),
		})
	}
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Generate predictions based on current standings
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction: " + err.Error(),
		})
	}
	
	// Assign probabilities based loosely on current standings
	totalProbability := 100.0
	var predictions []models.Prediction
	
	for i, stats := range standings {
		// Calculate remaining matches
		remainingMatches := 6 - stats.Played // Total of 6 matches per team
		maxPossiblePoints := stats.Points + (remainingMatches * 3)
		
		var probability float64
		if i == 0 {
			// Top team gets highest probability
			probability = 40.0 + (rand.Float64() * 10.0) // 40-50%
		} else if i == 1 {
			// Second team
			probability = 25.0 + (rand.Float64() * 10.0) // 25-35%
		} else if i == 2 {
			// Third team
			probability = 10.0 + (rand.Float64() * 10.0) // 10-20%
		} else {
			// Last team gets remaining probability
			probability = totalProbability - (predictions[0].PredictionPercentage + 
											predictions[1].PredictionPercentage + 
											predictions[2].PredictionPercentage)
		}
		
		// Ensure we don't go below 0 or have rounding issues
		if probability < 0 {
			probability = 0.1
		}
		
		prediction := models.Prediction{
			TeamID:             stats.Team.ID,
			Team:               stats.Team,
			PredictedPosition:  i + 1, // Current position
			PredictedPoints:    maxPossiblePoints,
			PredictionPercentage: probability,
		}
		
		predictions = append(predictions, prediction)
		
		// Insert prediction
		_, err := tx.Exec(
			"INSERT INTO predictions (team_id, predicted_position, predicted_points, prediction_percentage) VALUES ($1, $2, $3, $4)",
			prediction.TeamID, prediction.PredictedPosition, prediction.PredictedPoints, prediction.PredictionPercentage,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to insert prediction: " + err.Error(),
			})
		}
	}
	
	err = tx.Commit()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction: " + err.Error(),
		})
	}
	
	// Return the newly generated predictions
	return c.JSON(fiber.Map{
		"message": "Championship probabilities generated successfully",
		"predictions": predictions,
	})
} 