
package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"booking/internal/dto"
	"booking/internal/errs"

	"github.com/gin-gonic/gin"
)

type AvailabilityService interface {
	Get(ctx context.Context, roomID int64, from, to time.Time) (dto.AvailabilityResponse, error)
}

type AvailabilityHandler struct {
	svc AvailabilityService
}

func NewAvailabilityHandler(svc AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{svc: svc}
}

// Get godoc
// @Summary Проверить доступность комнаты
// @Description Возвращает информацию о занятых интервалах и доступности комнаты в указанном диапазоне
// @Tags availability
// @Produce json
// @Param id path int true "ID комнаты"
// @Param from query string true "Начало периода в RFC3339"
// @Param to query string true "Конец периода в RFC3339"
// @Success 200 {object} dto.AvailabilityResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /rooms/{id}/availability [get]
func (h *AvailabilityHandler) Get(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || roomID <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid room id"))
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "from/to are required"))
		return
	}

	from, err1 := time.Parse(time.RFC3339, fromStr)
	to, err2 := time.Parse(time.RFC3339, toStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "from/to must be RFC3339"))
		return
	}

	out, err := h.svc.Get(c.Request.Context(), roomID, from, to)
	if err != nil {
		switch err {
		case errs.ErrValidation:
			c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid interval"))
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		}
		return
	}

	c.JSON(http.StatusOK, out)
}