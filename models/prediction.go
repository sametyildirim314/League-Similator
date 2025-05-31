package models

import "time"

// Prediction represents a prediction for a team's final position
type Prediction struct {
	ID                 int       `json:"id"`
	TeamID             int       `json:"team_id"`
	Team               Team      `json:"team"`
	PredictedPosition  int       `json:"predicted_position"`
	PredictedPoints    int       `json:"predicted_points"`
	PredictionPercentage float64   `json:"prediction_percentage"`
	CreatedAt          time.Time `json:"created_at"`
} 