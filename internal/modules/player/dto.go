package player

import "github.com/yourname/football-api/internal/domain"

type CreatePlayerRequest struct {
	TeamID       uint                   `json:"team_id" binding:"required"`
	Name         string                 `json:"name" binding:"required,min=2,max=100"`
	Height       float32                `json:"height" binding:"required,gt=0"`
	Weight       float32                `json:"weight" binding:"required,gt=0"`
	Position     domain.PlayerPosition  `json:"position" binding:"required"`
	JerseyNumber int                    `json:"jersey_number" binding:"required,min=1,max=99"`
}

type UpdatePlayerRequest struct {
	Name         string                 `json:"name" binding:"omitempty,min=2,max=100"`
	Height       float32                `json:"height" binding:"omitempty,gt=0"`
	Weight       float32                `json:"weight" binding:"omitempty,gt=0"`
	Position     domain.PlayerPosition  `json:"position"`
	JerseyNumber int                    `json:"jersey_number" binding:"omitempty,min=1,max=99"`
	TeamID       uint                   `json:"team_id"`
}
