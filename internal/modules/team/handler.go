package team

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
	teams := r.Group("/teams")

	// Public (read)
	teams.GET("", h.GetAll)
	teams.GET("/:id", h.GetByID)

	// Protected (write) — butuh X-API-Key
	teams.Use(authMiddleware)
	teams.POST("", h.Create)
	teams.PUT("/:id", h.Update)
	teams.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary      Tambah tim baru
// @Tags         Teams
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body team.CreateTeamRequest true "Data tim"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /teams [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	team, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, ErrDuplicateName) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.InternalError(c, "gagal membuat tim")
		return
	}

	response.Created(c, "tim berhasil dibuat", team)
}

// GetAll godoc
// @Summary      Daftar semua tim
// @Tags         Teams
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /teams [get]
func (h *Handler) GetAll(c *gin.Context) {
	teams, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, "gagal mengambil data tim")
		return
	}
	response.OK(c, "data tim berhasil diambil", teams)
}

// GetByID godoc
// @Summary      Detail tim
// @Tags         Teams
// @Produce      json
// @Param        id path int true "ID Tim"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /teams/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	team, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrTeamNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal mengambil data tim")
		return
	}

	response.OK(c, "data tim berhasil diambil", team)
}

// Update godoc
// @Summary      Update tim
// @Tags         Teams
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Tim"
// @Param        body body team.UpdateTeamRequest true "Data tim yang diubah"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /teams/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	team, err := h.svc.Update(id, &req)
	if err != nil {
		if errors.Is(err, ErrTeamNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, ErrDuplicateName) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.InternalError(c, "gagal mengupdate tim")
		return
	}

	response.OK(c, "tim berhasil diupdate", team)
}

// Delete godoc
// @Summary      Hapus tim (soft delete)
// @Tags         Teams
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id path int true "ID Tim"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /teams/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrTeamNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal menghapus tim")
		return
	}

	response.OK(c, "tim berhasil dihapus", nil)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
