package match

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
	matches := r.Group("/matches")

	matches.GET("", h.GetAll)
	matches.GET("/:id", h.GetByID)

	matches.Use(authMiddleware)
	matches.POST("", h.Create)
	matches.PUT("/:id", h.Update)
	matches.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary      Buat jadwal pertandingan
// @Tags         Matches
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body match.CreateMatchRequest true "Data pertandingan"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /matches [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	match, err := h.svc.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSameTeam):
			response.BadRequest(c, err.Error())
		case errors.Is(err, ErrTeamNotFound):
			response.NotFound(c, err.Error())
		default:
			response.InternalError(c, "gagal membuat jadwal pertandingan")
		}
		return
	}

	response.Created(c, "jadwal pertandingan berhasil dibuat", match)
}

// GetAll godoc
// @Summary      Daftar semua pertandingan
// @Tags         Matches
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /matches [get]
func (h *Handler) GetAll(c *gin.Context) {
	matches, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, "gagal mengambil data pertandingan")
		return
	}
	response.OK(c, "data pertandingan berhasil diambil", matches)
}

// GetByID godoc
// @Summary      Detail pertandingan
// @Tags         Matches
// @Produce      json
// @Param        id path int true "ID Pertandingan"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /matches/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	match, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrMatchNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal mengambil data pertandingan")
		return
	}

	response.OK(c, "data pertandingan berhasil diambil", match)
}

// Update godoc
// @Summary      Update jadwal pertandingan
// @Tags         Matches
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Pertandingan"
// @Param        body body match.UpdateMatchRequest true "Data pertandingan yang diubah"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /matches/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req UpdateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	match, err := h.svc.Update(id, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMatchNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, ErrSameTeam):
			response.BadRequest(c, err.Error())
		case errors.Is(err, ErrTeamNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, ErrAlreadyCompleted):
			response.Error(c, http.StatusUnprocessableEntity, err.Error())
		default:
			response.InternalError(c, "gagal mengupdate pertandingan")
		}
		return
	}

	response.OK(c, "pertandingan berhasil diupdate", match)
}

// Delete godoc
// @Summary      Hapus jadwal pertandingan (soft delete)
// @Tags         Matches
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id path int true "ID Pertandingan"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /matches/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrMatchNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal menghapus pertandingan")
		return
	}

	response.OK(c, "pertandingan berhasil dihapus", nil)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
