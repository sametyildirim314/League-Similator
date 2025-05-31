package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/database"
	"github.com/sametyildirim314/insider_case/models"
)

// GetAllMatches handles the request to get all matches
func GetAllMatches(c *fiber.Ctx) error {
	// Query all matches directly from database
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, 
		       m.week, m.played, m.created_at,
		       ht.id, ht.name, at.id, at.name
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		ORDER BY m.week, m.id
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get matches: " + err.Error(),
		})
	}
	defer rows.Close()
	
	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var createdAt sql.NullTime
		
		err := rows.Scan(
			&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.HomeScore, &match.AwayScore,
			&match.Week, &match.Played, &createdAt,
			&match.HomeTeam.ID, &match.HomeTeam.Name, &match.AwayTeam.ID, &match.AwayTeam.Name,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan match: " + err.Error(),
			})
		}
		
		if createdAt.Valid {
			match.CreatedAt = createdAt.Time
		}
		
		matches = append(matches, match)
	}
	
	return c.JSON(matches)
}

// GetMatchesByWeek handles the request to get matches for a specific week
func GetMatchesByWeek(c *fiber.Ctx) error {
	week, err := strconv.Atoi(c.Params("week"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}
	
	// Query matches for the specific week directly from database
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, 
		       m.week, m.played, m.created_at,
		       ht.id, ht.name, at.id, at.name
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.week = $1
		ORDER BY m.id
	`
	
	rows, err := database.DB.Query(query, week)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get matches for week " + strconv.Itoa(week) + ": " + err.Error(),
		})
	}
	defer rows.Close()
	
	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var createdAt sql.NullTime
		
		err := rows.Scan(
			&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.HomeScore, &match.AwayScore,
			&match.Week, &match.Played, &createdAt,
			&match.HomeTeam.ID, &match.HomeTeam.Name, &match.AwayTeam.ID, &match.AwayTeam.Name,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan match: " + err.Error(),
			})
		}
		
		if createdAt.Valid {
			match.CreatedAt = createdAt.Time
		}
		
		matches = append(matches, match)
	}
	
	return c.JSON(matches)
}

// SimulateWeek handles the request to simulate matches for a specific week
func SimulateWeek(c *fiber.Ctx) error {
	week, err := strconv.Atoi(c.Params("week"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}
	
	// Check if fixtures exist
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM matches").Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check fixtures: " + err.Error(),
		})
	}
	
	// If no fixtures exist, generate them
	if count == 0 {
		if err := generateFixtures(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate fixtures: " + err.Error(),
			})
		}
	}
	
	// Check if previous weeks are simulated
	if week > 1 {
		var previousWeekCount int
		err = database.DB.QueryRow(
			"SELECT COUNT(*) FROM matches WHERE week < $1 AND played = false", 
			week,
		).Scan(&previousWeekCount)
		
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check previous weeks: " + err.Error(),
			})
		}
		
		if previousWeekCount > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to simulate week " + strconv.Itoa(week) + ": cannot simulate week " + strconv.Itoa(week) + ": previous weeks must be simulated first",
			})
		}
	}
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Get matches for the week
	rows, err := database.DB.Query(
		"SELECT id, home_team_id, away_team_id FROM matches WHERE week = $1 AND played = false",
		week,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get matches for week " + strconv.Itoa(week) + ": " + err.Error(),
		})
	}
	defer rows.Close()
	
	// For each match, simulate the result
	for rows.Next() {
		var matchID, homeTeamID, awayTeamID int
		err := rows.Scan(&matchID, &homeTeamID, &awayTeamID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan match: " + err.Error(),
			})
		}
		
		// Generate random scores (simple simulation)
		homeScore := rand.Intn(6) // 0-5 goals
		awayScore := rand.Intn(6) // 0-5 goals
		
		// Start a transaction
		tx, err := database.DB.Begin()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to start transaction: " + err.Error(),
			})
		}
		
		// Update match with scores and mark as played
		_, err = tx.Exec(
			"UPDATE matches SET home_score = $1, away_score = $2, played = true WHERE id = $3",
			homeScore, awayScore, matchID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update match: " + err.Error(),
			})
		}
		
		// Update league table for home team
		var homePoints int
		var homeWins, homeDraws, homeLosses int
		if homeScore > awayScore {
			// Home win
			homePoints = 3
			homeWins = 1
		} else if homeScore == awayScore {
			// Draw
			homePoints = 1
			homeDraws = 1
		} else {
			// Home loss
			homeLosses = 1
		}
		
		_, err = tx.Exec(`
			UPDATE league_table SET 
			points = points + $1,
			played = played + 1,
			wins = wins + $2,
			draws = draws + $3,
			losses = losses + $4,
			goals_for = goals_for + $5,
			goals_against = goals_against + $6,
			goal_difference = goal_difference + $7
			WHERE team_id = $8
		`,
			homePoints, homeWins, homeDraws, homeLosses, 
			homeScore, awayScore, homeScore-awayScore,
			homeTeamID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update home team stats: " + err.Error(),
			})
		}
		
		// Update league table for away team
		var awayPoints int
		var awayWins, awayDraws, awayLosses int
		if awayScore > homeScore {
			// Away win
			awayPoints = 3
			awayWins = 1
		} else if awayScore == homeScore {
			// Draw
			awayPoints = 1
			awayDraws = 1
		} else {
			// Away loss
			awayLosses = 1
		}
		
		_, err = tx.Exec(`
			UPDATE league_table SET 
			points = points + $1,
			played = played + 1,
			wins = wins + $2,
			draws = draws + $3,
			losses = losses + $4,
			goals_for = goals_for + $5,
			goals_against = goals_against + $6,
			goal_difference = goal_difference + $7
			WHERE team_id = $8
		`,
			awayPoints, awayWins, awayDraws, awayLosses, 
			awayScore, homeScore, awayScore-homeScore,
			awayTeamID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update away team stats: " + err.Error(),
			})
		}
		
		// Commit transaction
		err = tx.Commit()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to commit transaction: " + err.Error(),
			})
		}
	}
	
	// Get updated matches for the week
	updatedMatches, err := getMatchesByWeek(week)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get matches after simulation",
		})
	}
	
	// After week 4, automatically generate championship predictions
	if week == 4 {
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
			// Just log the error but don't fail the whole request
			fmt.Printf("Failed to get league standings for predictions: %v\n", err)
		} else {
			var standings []models.TeamStats
			for rows.Next() {
				var stats models.TeamStats
				err := rows.Scan(
					&stats.Points, &stats.Played, &stats.Wins, &stats.Draws, &stats.Losses,
					&stats.GoalsFor, &stats.GoalsAgainst, &stats.GoalDifference,
					&stats.Team.ID, &stats.Team.Name,
				)
				if err != nil {
					fmt.Printf("Failed to scan team stats: %v\n", err)
					continue
				}
				
				standings = append(standings, stats)
			}
			rows.Close()
			
			// Clear existing predictions
			_, err = database.DB.Exec("DELETE FROM predictions")
			if err != nil {
				fmt.Printf("Failed to clear existing predictions: %v\n", err)
			} else {
				// Seed random number generator
				rand.Seed(time.Now().UnixNano())
				
				// Generate predictions based on current standings
				tx, err := database.DB.Begin()
				if err != nil {
					fmt.Printf("Failed to start transaction for predictions: %v\n", err)
				} else {
					// Assign probabilities based loosely on current standings
					totalProbability := 100.0
					var generatedPredictions []models.Prediction
					
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
							probability = totalProbability
							if len(generatedPredictions) >= 3 {
								probability = totalProbability - (generatedPredictions[0].PredictionPercentage + 
														generatedPredictions[1].PredictionPercentage + 
														generatedPredictions[2].PredictionPercentage)
							}
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
						
						generatedPredictions = append(generatedPredictions, prediction)
						
						// Insert prediction
						_, err := tx.Exec(
							"INSERT INTO predictions (team_id, predicted_position, predicted_points, prediction_percentage) VALUES ($1, $2, $3, $4)",
							prediction.TeamID, prediction.PredictedPosition, prediction.PredictedPoints, prediction.PredictionPercentage,
						)
						if err != nil {
							tx.Rollback()
							fmt.Printf("Failed to insert prediction: %v\n", err)
							break
						}
					}
					
					err = tx.Commit()
					if err != nil {
						fmt.Printf("Failed to commit transaction for predictions: %v\n", err)
					} else {
						return c.JSON(fiber.Map{
							"message": "Successfully simulated week " + strconv.Itoa(week) + " and generated championship predictions",
							"matches": updatedMatches,
							"predictions": generatedPredictions,
						})
					}
				}
			}
		}
	} else if week >= 4 {
		// Week 4 and above should update predictions too
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
			// Just log the error but don't fail the whole request
			fmt.Printf("Failed to get league standings for predictions: %v\n", err)
		} else {
			var standings []models.TeamStats
			for rows.Next() {
				var stats models.TeamStats
				err := rows.Scan(
					&stats.Points, &stats.Played, &stats.Wins, &stats.Draws, &stats.Losses,
					&stats.GoalsFor, &stats.GoalsAgainst, &stats.GoalDifference,
					&stats.Team.ID, &stats.Team.Name,
				)
				if err != nil {
					fmt.Printf("Failed to scan team stats: %v\n", err)
					continue
				}
				
				standings = append(standings, stats)
			}
			rows.Close()
			
			// Clear existing predictions
			_, err = database.DB.Exec("DELETE FROM predictions")
			if err != nil {
				fmt.Printf("Failed to clear existing predictions: %v\n", err)
			} else {
				// Seed random number generator
				rand.Seed(time.Now().UnixNano())
				
				// Generate predictions based on current standings
				tx, err := database.DB.Begin()
				if err != nil {
					fmt.Printf("Failed to start transaction for predictions: %v\n", err)
				} else {
					// Assign probabilities based loosely on current standings
					totalProbability := 100.0
					var generatedPredictions []models.Prediction
					
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
							probability = totalProbability
							if len(generatedPredictions) >= 3 {
								probability = totalProbability - (generatedPredictions[0].PredictionPercentage + 
														generatedPredictions[1].PredictionPercentage + 
														generatedPredictions[2].PredictionPercentage)
							}
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
						
						generatedPredictions = append(generatedPredictions, prediction)
						
						// Insert prediction
						_, err := tx.Exec(
							"INSERT INTO predictions (team_id, predicted_position, predicted_points, prediction_percentage) VALUES ($1, $2, $3, $4)",
							prediction.TeamID, prediction.PredictedPosition, prediction.PredictedPoints, prediction.PredictionPercentage,
						)
						if err != nil {
							tx.Rollback()
							fmt.Printf("Failed to insert prediction: %v\n", err)
							break
						}
					}
					
					err = tx.Commit()
					if err != nil {
						fmt.Printf("Failed to commit transaction for predictions: %v\n", err)
					} else {
						return c.JSON(fiber.Map{
							"message": "Successfully simulated week " + strconv.Itoa(week) + " and generated championship predictions",
							"matches": updatedMatches,
							"predictions": generatedPredictions,
						})
					}
				}
			}
		}
	}
	
	return c.JSON(fiber.Map{
		"message": "Successfully simulated week " + strconv.Itoa(week),
		"matches": updatedMatches,
	})
}

// SimulateAllRemainingMatches handles the request to simulate all remaining matches
func SimulateAllRemainingMatches(c *fiber.Ctx) error {
	// Check if fixtures exist
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM matches").Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check fixtures: " + err.Error(),
		})
	}
	
	// If no fixtures exist, generate them
	if count == 0 {
		if err := generateFixtures(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate fixtures: " + err.Error(),
			})
		}
	}
	
	// Get all unplayed matches ordered by week
	rows, err := database.DB.Query(
		"SELECT id, home_team_id, away_team_id, week FROM matches WHERE played = false ORDER BY week, id",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unplayed matches: " + err.Error(),
		})
	}
	defer rows.Close()
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// For each match, simulate the result
	for rows.Next() {
		var matchID, homeTeamID, awayTeamID, week int
		err := rows.Scan(&matchID, &homeTeamID, &awayTeamID, &week)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan match: " + err.Error(),
			})
		}
		
		// Generate random scores (simple simulation)
		homeScore := rand.Intn(6) // 0-5 goals
		awayScore := rand.Intn(6) // 0-5 goals
		
		// Start a transaction
		tx, err := database.DB.Begin()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to start transaction: " + err.Error(),
			})
		}
		
		// Update match with scores and mark as played
		_, err = tx.Exec(
			"UPDATE matches SET home_score = $1, away_score = $2, played = true WHERE id = $3",
			homeScore, awayScore, matchID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update match: " + err.Error(),
			})
		}
		
		// Update league table for home team
		var homePoints int
		var homeWins, homeDraws, homeLosses int
		if homeScore > awayScore {
			// Home win
			homePoints = 3
			homeWins = 1
		} else if homeScore == awayScore {
			// Draw
			homePoints = 1
			homeDraws = 1
		} else {
			// Home loss
			homeLosses = 1
		}
		
		_, err = tx.Exec(`
			UPDATE league_table SET 
			points = points + $1,
			played = played + 1,
			wins = wins + $2,
			draws = draws + $3,
			losses = losses + $4,
			goals_for = goals_for + $5,
			goals_against = goals_against + $6,
			goal_difference = goal_difference + $7
			WHERE team_id = $8
		`,
			homePoints, homeWins, homeDraws, homeLosses, 
			homeScore, awayScore, homeScore-awayScore,
			homeTeamID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update home team stats: " + err.Error(),
			})
		}
		
		// Update league table for away team
		var awayPoints int
		var awayWins, awayDraws, awayLosses int
		if awayScore > homeScore {
			// Away win
			awayPoints = 3
			awayWins = 1
		} else if awayScore == homeScore {
			// Draw
			awayPoints = 1
			awayDraws = 1
		} else {
			// Away loss
			awayLosses = 1
		}
		
		_, err = tx.Exec(`
			UPDATE league_table SET 
			points = points + $1,
			played = played + 1,
			wins = wins + $2,
			draws = draws + $3,
			losses = losses + $4,
			goals_for = goals_for + $5,
			goals_against = goals_against + $6,
			goal_difference = goal_difference + $7
			WHERE team_id = $8
		`,
			awayPoints, awayWins, awayDraws, awayLosses, 
			awayScore, homeScore, awayScore-homeScore,
			awayTeamID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update away team stats: " + err.Error(),
			})
		}
		
		// Commit transaction
		err = tx.Commit()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to commit transaction: " + err.Error(),
			})
		}
	}
	
	// Get all matches
	allMatches, err := getAllMatches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get matches after simulation",
		})
	}

	// Generate championship predictions after all matches are simulated
	// Get current league standings from database
	standingsQuery := `
		SELECT lt.points, lt.played, lt.wins, lt.draws, lt.losses, 
			lt.goals_for, lt.goals_against, lt.goal_difference,
			t.id, t.name
		FROM league_table lt
		JOIN teams t ON lt.team_id = t.id
		ORDER BY lt.points DESC, lt.goal_difference DESC, lt.goals_for DESC
	`
	
	rows, err = database.DB.Query(standingsQuery)
	if err != nil {
		// Just log the error but don't fail the whole request
		fmt.Printf("Failed to get league standings for predictions: %v\n", err)
		return c.JSON(fiber.Map{
			"message": "Successfully simulated all remaining matches",
			"matches": allMatches,
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
			fmt.Printf("Failed to scan team stats: %v\n", err)
			continue
		}
		
		standings = append(standings, stats)
	}
	rows.Close()
	
	// Clear existing predictions
	_, err = database.DB.Exec("DELETE FROM predictions")
	if err != nil {
		fmt.Printf("Failed to clear existing predictions: %v\n", err)
		return c.JSON(fiber.Map{
			"message": "Successfully simulated all remaining matches",
			"matches": allMatches,
		})
	}
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Generate predictions based on current standings
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Printf("Failed to start transaction for predictions: %v\n", err)
		return c.JSON(fiber.Map{
			"message": "Successfully simulated all remaining matches",
			"matches": allMatches,
		})
	}
	
	// Assign probabilities based loosely on current standings
	totalProbability := 100.0
	var generatedPredictions []models.Prediction
	
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
			probability = totalProbability
			if len(generatedPredictions) >= 3 {
				probability = totalProbability - (generatedPredictions[0].PredictionPercentage + 
								generatedPredictions[1].PredictionPercentage + 
								generatedPredictions[2].PredictionPercentage)
			}
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
		
		generatedPredictions = append(generatedPredictions, prediction)
		
		// Insert prediction
		_, err := tx.Exec(
			"INSERT INTO predictions (team_id, predicted_position, predicted_points, prediction_percentage) VALUES ($1, $2, $3, $4)",
			prediction.TeamID, prediction.PredictedPosition, prediction.PredictedPoints, prediction.PredictionPercentage,
		)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Failed to insert prediction: %v\n", err)
			return c.JSON(fiber.Map{
				"message": "Successfully simulated all remaining matches",
				"matches": allMatches,
			})
		}
	}
	
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Failed to commit transaction for predictions: %v\n", err)
		return c.JSON(fiber.Map{
			"message": "Successfully simulated all remaining matches",
			"matches": allMatches,
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Successfully simulated all remaining matches and generated championship predictions",
		"matches": allMatches,
		"predictions": generatedPredictions,
	})
}

// generateFixtures generates fixtures for the league
func generateFixtures() error {
	// Get all teams
	rows, err := database.DB.Query("SELECT id FROM teams")
	if err != nil {
		return err
	}
	defer rows.Close()
	
	var teams []int
	for rows.Next() {
		var teamID int
		err := rows.Scan(&teamID)
		if err != nil {
			return err
		}
		teams = append(teams, teamID)
	}
	
	// Check if we have enough teams
	if len(teams) < 2 {
		return errors.New("not enough teams to generate fixtures")
	}
	
	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// For our specific case with 4 teams, we'll manually create a balanced fixture
	// that spreads across 6 weeks
	if len(teams) == 4 {
		// First 3 weeks - each team plays against others once
		fixtures := []struct {
			home int
			away int
			week int
		}{
			// Week 1
			{teams[0], teams[1], 1},
			{teams[2], teams[3], 1},
			
			// Week 2
			{teams[0], teams[2], 2},
			{teams[1], teams[3], 2},
			
			// Week 3
			{teams[0], teams[3], 3},
			{teams[1], teams[2], 3},
			
			// Week 4 - reverse fixtures
			{teams[1], teams[0], 4},
			{teams[3], teams[2], 4},
			
			// Week 5
			{teams[2], teams[0], 5},
			{teams[3], teams[1], 5},
			
			// Week 6
			{teams[3], teams[0], 6},
			{teams[2], teams[1], 6},
		}
		
		for _, fixture := range fixtures {
			_, err = tx.Exec(
				"INSERT INTO matches (home_team_id, away_team_id, week, played) VALUES ($1, $2, $3, false)",
				fixture.home, fixture.away, fixture.week,
			)
			if err != nil {
				return err
			}
		}
		
		return tx.Commit()
	}
	
	// For other number of teams, use the algorithm
	totalWeeks := 6 // Fix to 6 weeks as shown in the example
	
	// Create a list of all possible matchups
	var matchups []struct {
		homeTeam int
		awayTeam int
	}
	
	for i := 0; i < len(teams); i++ {
		for j := 0; j < len(teams); j++ {
			if i != j { // Teams don't play against themselves
				matchups = append(matchups, struct {
					homeTeam int
					awayTeam int
				}{
					homeTeam: teams[i],
					awayTeam: teams[j],
				})
			}
		}
	}
	
	// Shuffle matchups to randomize the fixture
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(matchups), func(i, j int) {
		matchups[i], matchups[j] = matchups[j], matchups[i]
	})
	
	// Distribute matchups across weeks
	for i, matchup := range matchups {
		week := (i % totalWeeks) + 1 // Weeks are 1-indexed
		
		// Insert fixture
		_, err = tx.Exec(
			"INSERT INTO matches (home_team_id, away_team_id, week, played) VALUES ($1, $2, $3, false)",
			matchup.homeTeam, matchup.awayTeam, week,
		)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// getAllMatches returns all matches from the database
func getAllMatches() ([]models.Match, error) {
	// Query all matches directly from database
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, 
		       m.week, m.played, m.created_at,
		       ht.id, ht.name, at.id, at.name
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		ORDER BY m.week, m.id
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %v", err)
	}
	defer rows.Close()
	
	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var createdAt sql.NullTime
		
		err := rows.Scan(
			&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.HomeScore, &match.AwayScore,
			&match.Week, &match.Played, &createdAt,
			&match.HomeTeam.ID, &match.HomeTeam.Name, &match.AwayTeam.ID, &match.AwayTeam.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %v", err)
		}
		
		if createdAt.Valid {
			match.CreatedAt = createdAt.Time
		}
		
		matches = append(matches, match)
	}
	
	return matches, nil
}

// getMatchesByWeek returns matches for a specific week
func getMatchesByWeek(week int) ([]models.Match, error) {
	// Query matches for the specific week directly from database
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, 
		       m.week, m.played, m.created_at,
		       ht.id, ht.name, at.id, at.name
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.week = $1
		ORDER BY m.id
	`
	
	rows, err := database.DB.Query(query, week)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches for week %d: %v", week, err)
	}
	defer rows.Close()
	
	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var createdAt sql.NullTime
		
		err := rows.Scan(
			&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.HomeScore, &match.AwayScore,
			&match.Week, &match.Played, &createdAt,
			&match.HomeTeam.ID, &match.HomeTeam.Name, &match.AwayTeam.ID, &match.AwayTeam.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %v", err)
		}
		
		if createdAt.Valid {
			match.CreatedAt = createdAt.Time
		}
		
		matches = append(matches, match)
	}
	
	return matches, nil
} 