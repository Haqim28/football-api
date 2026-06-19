package goal

// SubmitResultRequest dipakai untuk submit hasil pertandingan beserta daftar gol
type SubmitResultRequest struct {
	MatchID   uint          `json:"match_id" binding:"required"`
	HomeScore int           `json:"home_score" binding:"min=0"`
	AwayScore int           `json:"away_score" binding:"min=0"`
	Goals     []GoalRequest `json:"goals" binding:"required"`
}

type GoalRequest struct {
	PlayerID uint `json:"player_id" binding:"required"`
	Minute   int  `json:"minute" binding:"required,min=1,max=120"`
}
