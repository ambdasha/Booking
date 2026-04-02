package handlers

import (
	"context"
	"net/http"
	"strconv"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"
	"booking/internal/httpx/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReservationService interface {
	Create(ctx context.Context, userID int64, req dto.CreateReservationRequest) (domain.Reservation, error)
	Get(ctx context.Context, actorID int64, actorRole string, id int64) (domain.Reservation, error)
	ListMy(ctx context.Context, userID int64, status string) ([]domain.Reservation, error)
	Cancel(ctx context.Context, actorID int64, actorRole string, id int64, reason string) error
}

type ReservationsHandler struct {
	svc       ReservationService
	validator *validator.Validate
}

func NewReservationsHandler(svc ReservationService) *ReservationsHandler {
	return &ReservationsHandler{svc: svc, validator: validator.New()}
}

// Create godoc
// @Summary Создать бронирование
// @Description Создает новое бронирование комнаты
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateReservationRequest true "Данные бронирования"
// @Success 201 {object} dto.ReservationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /reservations [post]
func (h *ReservationsHandler) Create(c *gin.Context) {
	uid, _ := middleware.MustUser(c)

	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	// validator не всегда идеально ловит zero time, поэтому проверяем руками в сервисе
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	res, err := h.svc.Create(c.Request.Context(), uid, req)
	if err != nil {
		switch err {
		case errs.ErrValidation:
			c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid interval"))
			return
		case errs.ErrConflict:
			c.JSON(http.StatusConflict, errorResp("reservation_conflict", "selected time slot is already booked"))
			return
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, errorResp("not_found", "room not found"))
			return
		case errs.ErrForbidden:
			c.JSON(http.StatusForbidden, errorResp("forbidden", "room is inactive"))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
			return
		}
	}

	c.JSON(http.StatusCreated, dto.ReservationResponse{
		ID:        res.ID,
		UserID:    res.UserID,
		RoomID:    res.RoomID,
		Status:    res.Status,
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		CreatedAt: res.CreatedAt,
	})
}

// MyList godoc
// @Summary Мои бронирования
// @Description Возвращает список бронирований текущего пользователя
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ReservationResponse
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/my [get]
func (h *ReservationsHandler) MyList(c *gin.Context) {
	uid, _ := middleware.MustUser(c)
	status := c.Query("status") // можно '', confirmed, cancelled...

	items, err := h.svc.ListMy(c.Request.Context(), uid, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	out := make([]dto.ReservationResponse, 0, len(items))
	for _, res := range items {
		out = append(out, dto.ReservationResponse{
			ID:        res.ID,
			UserID:    res.UserID,
			RoomID:    res.RoomID,
			Status:    res.Status,
			StartTime: res.StartTime,
			EndTime:   res.EndTime,
			CreatedAt: res.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, out)
}

// Get godoc
// @Summary Получить бронирование
// @Description Возвращает бронирование по id
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID бронирования"
// @Success 200 {object} dto.ReservationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reservations/{id} [get]
func (h *ReservationsHandler) Get(c *gin.Context) {
	uid, role := middleware.MustUser(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		return
	}

	res, err := h.svc.Get(c.Request.Context(), uid, role, id)
	if err != nil {
		switch err {
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, errorResp("not_found", "reservation not found"))
			return
		case errs.ErrForbidden:
			c.JSON(http.StatusForbidden, errorResp("forbidden", "insufficient permissions"))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
			return
		}
	}

	c.JSON(http.StatusOK, dto.ReservationResponse{
		ID:        res.ID,
		UserID:    res.UserID,
		RoomID:    res.RoomID,
		Status:    res.Status,
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		CreatedAt: res.CreatedAt,
	})
}

// Cancel godoc
// @Summary Отменить бронирование
// @Description Отменяет бронирование по id
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID бронирования"
// @Success 200 {object} dto.ReservationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [post]
func (h *ReservationsHandler) Cancel(c *gin.Context) {
	uid, role := middleware.MustUser(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		return
	}

	var req dto.CancelReservationRequest
	_ = c.ShouldBindJSON(&req) 
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	err = h.svc.Cancel(c.Request.Context(), uid, role, id, req.Reason)
	if err != nil {
		switch err {
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, errorResp("not_found", "reservation not found"))
			return
		case errs.ErrForbidden:
			c.JSON(http.StatusForbidden, errorResp("forbidden", "insufficient permissions"))
			return
		case errs.ErrConflict:
			c.JSON(http.StatusConflict, errorResp("conflict", "cannot cancel this reservation"))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
			return
		}
	}

	c.Status(http.StatusNoContent)
}