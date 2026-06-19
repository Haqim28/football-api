package team

// ─── Request DTOs ─────────────────────────────────────────────────────────────

type CreateTeamRequest struct {
	Name                string `json:"name" binding:"required,min=2,max=100"`
	Logo                string `json:"logo"`
	Founded             int    `json:"founded" binding:"required,min=1800,max=2100"`
	HeadquartersAddress string `json:"headquarters_address" binding:"required"`
	HeadquartersCity    string `json:"headquarters_city" binding:"required"`
}

type UpdateTeamRequest struct {
	Name                string `json:"name" binding:"omitempty,min=2,max=100"`
	Logo                string `json:"logo"`
	Founded             int    `json:"founded" binding:"omitempty,min=1800,max=2100"`
	HeadquartersAddress string `json:"headquarters_address"`
	HeadquartersCity    string `json:"headquarters_city"`
}
