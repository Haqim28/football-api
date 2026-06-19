package report

import (
	"errors"
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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	reports := r.Group("/reports")
	reports.GET("", h.GetAll)
	reports.GET("/matches/:id", h.GetByMatch)
}

// GetAll godoc
// @Summary      Laporan semua pertandingan
// @Tags         Reports
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /reports [get]
func (h *Handler) GetAll(c *gin.Context) {
	reports, err := h.svc.GetAllReports()
	if err != nil {
		response.InternalError(c, "gagal mengambil data laporan")
		return
	}
	response.OK(c, "laporan berhasil diambil", reports)
}

// GetByMatch godoc
// @Summary      Laporan detail per pertandingan
// @Description  Menampilkan jadwal, skor, status akhir (Tim Home Menang/Tim Away Menang/Draw), top scorer, dan akumulasi kemenangan home & away team.
// @Tags         Reports
// @Produce      json
// @Param        id path int true "ID Pertandingan"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /reports/matches/{id} [get]
func (h *Handler) GetByMatch(c *gin.Context) {
	rawID := c.Param("id")
	parsedID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	report, err := h.svc.GetMatchReport(uint(parsedID))
	if err != nil {
		if errors.Is(err, ErrMatchNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "gagal mengambil laporan")
		return
	}

	response.OK(c, "laporan berhasil diambil", report)
}
