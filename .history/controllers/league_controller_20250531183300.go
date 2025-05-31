package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sametyildirim314/insider_case/database"
	"github.com/sametyildirim314/insider_case/models"
)


func GetLeagueTable(c *fiber.Ctx) error {

	query := `
		SELECT lt.points, lt.played, lt.wins, lt.draws, lt.losses, 
		       lt.goals_for, lt.goals_against, lt.goal_difference,
		       t.id, t.name
		FROM league_table lt
		JOIN teams t ON lt.team_id = t.id
		ORDER BY lt.points DESC, lt.goal_difference DESC, lt.goals_for DESC
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get league table: " + err.Error(),
		})
	}
	defer rows.Close()
	
	var teamStats []models.TeamStats
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
		
		teamStats = append(teamStats, stats)
	}
	
	return c.JSON(teamStats)
}

func GetLeagueTableForWeek(c *fiber.Ctx) error {
	week, err := strconv.Atoi(c.Params("week"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}
	
	// Check if any matches have been played for this week
	var matchesPlayed int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM matches WHERE week <= $1 AND played = true", week).Scan(&matchesPlayed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check matches: " + err.Error(),
		})
	}
	
	if matchesPlayed == 0 {
		// No matches played yet, return empty table with teams
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
		
		var teamStats []models.TeamStats
		for _, team := range teams {
			stats := models.TeamStats{
				Team:          team,
				Points:        0,
				Played:        0,
				Wins:          0,
				Draws:         0,
				Losses:        0,
				GoalsFor:      0,
				GoalsAgainst:  0,
				GoalDifference: 0,
			}
			teamStats = append(teamStats, stats)
		}
		
		return c.JSON(teamStats)
	}
	
	// Otherwise, use the current table for simplicity
	return GetLeagueTable(c)
} 