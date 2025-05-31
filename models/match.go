package models

import "time"

// Match represents a match between two teams
type Match struct {
	ID          int       `json:"id"`
	HomeTeamID  int       `json:"home_team_id"`
	AwayTeamID  int       `json:"away_team_id"`
	HomeTeam    Team      `json:"home_team"`
	AwayTeam    Team      `json:"away_team"`
	HomeScore   *int      `json:"home_score"`
	AwayScore   *int      `json:"away_score"`
	Week        int       `json:"week"`
	Played      bool      `json:"played"`
	CreatedAt   time.Time `json:"created_at"`
} 