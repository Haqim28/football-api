package domain

import (
	"time"

	"gorm.io/gorm"
)

// ─── Team ────────────────────────────────────────────────────────────────────

type Team struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Logo         string         `gorm:"type:varchar(255)" json:"logo"`
	Founded      int            `gorm:"not null" json:"founded"`
	HeadquartersAddress string  `gorm:"type:varchar(255);not null" json:"headquarters_address"`
	HeadquartersCity    string  `gorm:"type:varchar(100);not null" json:"headquarters_city"`
	Players      []Player       `gorm:"foreignKey:TeamID" json:"players,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ─── Player ───────────────────────────────────────────────────────────────────

type PlayerPosition string

const (
	PositionForward    PlayerPosition = "penyerang"
	PositionMidfielder PlayerPosition = "gelandang"
	PositionDefender   PlayerPosition = "bertahan"
	PositionGoalkeeper PlayerPosition = "penjaga_gawang"
)

type Player struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TeamID       uint           `gorm:"not null;index" json:"team_id"`
	Team         *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Height       float32        `gorm:"not null" json:"height"` // cm
	Weight       float32        `gorm:"not null" json:"weight"` // kg
	Position     PlayerPosition `gorm:"type:varchar(30);not null" json:"position"`
	JerseyNumber int            `gorm:"not null" json:"jersey_number"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ─── Match ────────────────────────────────────────────────────────────────────

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusCompleted MatchStatus = "completed"
)

type Match struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MatchDate  time.Time      `gorm:"not null" json:"match_date"`
	HomeTeamID uint           `gorm:"not null;index" json:"home_team_id"`
	AwayTeamID uint           `gorm:"not null;index" json:"away_team_id"`
	HomeTeam   *Team          `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam   *Team          `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	HomeScore  int            `gorm:"default:0" json:"home_score"`
	AwayScore  int            `gorm:"default:0" json:"away_score"`
	Status     MatchStatus    `gorm:"type:varchar(20);default:'scheduled'" json:"status"`
	Goals      []Goal         `gorm:"foreignKey:MatchID" json:"goals,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ─── Goal ─────────────────────────────────────────────────────────────────────

type Goal struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MatchID   uint           `gorm:"not null;index" json:"match_id"`
	PlayerID  uint           `gorm:"not null;index" json:"player_id"`
	TeamID    uint           `gorm:"not null;index" json:"team_id"`
	Minute    int            `gorm:"not null" json:"minute"` // menit ke-berapa gol terjadi
	Match     *Match         `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	Player    *Player        `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Team      *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
