package models

// Team represents a football team in the Premier League
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TeamStats represents a team with its statistics for the league table
type TeamStats struct {
	Team          Team `json:"team"`
	Points        int  `json:"points"`
	Played        int  `json:"played"`
	Wins          int  `json:"wins"`
	Draws         int  `json:"draws"`
	Losses        int  `json:"losses"`
	GoalsFor      int  `json:"goals_for"`
	GoalsAgainst  int  `json:"goals_against"`
	GoalDifference int  `json:"goal_difference"`
} 