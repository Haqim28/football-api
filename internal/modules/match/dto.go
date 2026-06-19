package match

import "time"

type CreateMatchRequest struct {
	MatchDate  time.Time `json:"match_date" binding:"required"`
	HomeTeamID uint      `json:"home_team_id" binding:"required"`
	AwayTeamID uint      `json:"away_team_id" binding:"required"`
}

type UpdateMatchRequest struct {
	MatchDate  *time.Time `json:"match_date"`
	HomeTeamID uint       `json:"home_team_id"`
	AwayTeamID uint       `json:"away_team_id"`
}
