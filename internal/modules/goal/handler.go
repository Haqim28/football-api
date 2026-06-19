package goal

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
	// POST /api/v1/matches/:id/result — submit hasil pertandingan
	r.POST("/matches/:id/result", authMiddleware, h.SubmitResult)
}

// SubmitResult godoc
// @Summary      Submit hasil pertandingan
// @Description  Menyimpan skor akhir dan daftar pencetak gol. Jumlah gol pemain tim home harus sama dengan home_score, begitu juga away_score untuk tim away. Match otomatis berubah jadi 'completed' setelah submit.
// @Tags         Matches
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Pertandingan"
// @Param        body body goal.SubmitResultRequest true "Hasil pertandingan & daftar gol"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /matches/{id}/result [post]
func (h *Handler) SubmitResult(c *gin.Context) {
	// Ambil match_id dari URL param
	rawID := c.Param("id")
	parsedID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID pertandingan tidak valid")
		return
	}
	matchID := uint(parsedID)

	var req SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Paksa match_id dari URL (ignore body match_id kalau ada)
	req.MatchID = matchID

	goals, err := h.svc.SubmitResult(&req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMatchNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, ErrMatchAlreadyDone):
			response.Error(c, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, ErrPlayerNotInMatch):
			response.BadRequest(c, err.Error())
		case errors.Is(err, ErrScoreMismatch):
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, "gagal menyimpan hasil pertandingan")
		}
		return
	}

	response.OK(c, "hasil pertandingan berhasil disimpan", goals)
}
