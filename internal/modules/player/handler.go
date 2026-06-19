package player

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourname/football-api/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	players := r.Group("/players")

	players.GET("", h.GetAll)
	players.GET("/:id", h.GetByID)

	players.Use(authMiddleware)
	players.POST("", h.Create)
	players.PUT("/:id", h.Update)
	players.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary      Tambah pemain baru
// @Tags         Players
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body player.CreatePlayerRequest true "Data pemain"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /players [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	player, err := h.svc.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateJersey):
			response.Error(c, http.StatusConflict, err.Error())
		case errors.Is(err, ErrTeamNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, ErrInvalidPosition):
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, "gagal membuat pemain")
		}
		return
	}

	response.Created(c, "pemain berhasil ditambahkan", player)
}

// GetAll godoc
// @Summary      Daftar semua pemain
// @Tags         Players
// @Produce      json
// @Param        team_id query int false "Filter pemain berdasarkan tim"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /players [get]
func (h *Handler) GetAll(c *gin.Context) {
	var teamID uint
	if teamIDStr := c.Query("team_id"); teamIDStr != "" {
		id, err := strconv.ParseUint(teamIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "team_id tidak valid")
			return
		}
		teamID = uint(id)
	}

	players, err := h.svc.GetAll(teamID)
	if err != nil {
		response.InternalError(c, "gagal mengambil data pemain")
		return
	}
	response.OK(c, "data pemain berhasil diambil", players)
}

// GetByID godoc
// @Summary      Detail pemain
// @Tags         Players
// @Produce      json
// @Param        id path int true "ID Pemain"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /players/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	player, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrPlayerNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal mengambil data pemain")
		return
	}

	response.OK(c, "data pemain berhasil diambil", player)
}

// Update godoc
// @Summary      Update pemain
// @Tags         Players
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Pemain"
// @Param        body body player.UpdatePlayerRequest true "Data pemain yang diubah"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /players/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	player, err := h.svc.Update(id, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrPlayerNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, ErrDuplicateJersey):
			response.Error(c, http.StatusConflict, err.Error())
		case errors.Is(err, ErrInvalidPosition):
			response.BadRequest(c, err.Error())
		case errors.Is(err, ErrTeamNotFound):
			response.NotFound(c, err.Error())
		default:
			response.InternalError(c, "gagal mengupdate pemain")
		}
		return
	}

	response.OK(c, "pemain berhasil diupdate", player)
}

// Delete godoc
// @Summary      Hapus pemain (soft delete)
// @Tags         Players
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id path int true "ID Pemain"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /players/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrPlayerNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal menghapus pemain")
		return
	}

	response.OK(c, "pemain berhasil dihapus", nil)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
