package handlers
//HTTP+валидация+ответ
//парсит json,валидирует dto, дергает service/rooms.go, возвращает json

import (
	"context"
	"net/http"
	"strconv"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoomService interface {
	List(ctx context.Context) ([]domain.Room, error)
	Get(ctx context.Context, id int64) (domain.Room, error)
	Create(ctx context.Context, req dto.CreateRoomRequest) (domain.Room, error)
	Update(ctx context.Context, id int64, req dto.UpdateRoomRequest) (domain.Room, error)
	Deactivate(ctx context.Context, id int64) error
}

type RoomsHandler struct {
	svc       RoomService
	validator *validator.Validate
}

func NewRoomsHandler(svc RoomService) *RoomsHandler {
	return &RoomsHandler{
		svc:       svc,
		validator: validator.New(),
	}
}

// List godoc
// @Summary Список комнат
// @Description Возвращает список активных комнат
// @Tags rooms
// @Produce json
// @Success 200 {array} dto.RoomResponse
// @Failure 500 {object} map[string]interface{}
// @Router /rooms [get]
func (h *RoomsHandler) List(c *gin.Context) {
	rooms, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}
//контроль того что будет показано клиенту т.к сервис возвращает []domain.Room, наружу отдается []dto.RoomResponse
	out := make([]dto.RoomResponse, 0, len(rooms))
	for _, r := range rooms {
		out = append(out, roomToResponse(r))
	}

	c.JSON(http.StatusOK, out)
}


// Get godoc
// @Summary Получить комнату
// @Description Возвращает комнату по id
// @Tags rooms
// @Produce json
// @Param id path int true "ID комнаты"
// @Success 200 {object} dto.RoomResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /rooms/{id} [get]
func (h *RoomsHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		return
	}

	room, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == errs.ErrNotFound {
			c.JSON(http.StatusNotFound, errorResp("not_found", "room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	c.JSON(http.StatusOK, roomToResponse(room))
}

// Create godoc
// @Summary Создать комнату
// @Description Создает новую комнату. Доступно только администратору
// @Tags admin-rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateRoomRequest true "Данные комнаты"
// @Success 201 {object} dto.RoomResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/rooms [post]
func (h *RoomsHandler) Create(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	room, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	c.JSON(http.StatusCreated, roomToResponse(room))
}


// Update godoc
// @Summary Обновить комнату
// @Description Обновляет комнату по id. Доступно только администратору
// @Tags admin-rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комнаты"
// @Param input body dto.UpdateRoomRequest true "Новые данные комнаты"
// @Success 200 {object} dto.RoomResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/rooms/{id} [put]
func (h *RoomsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		return
	}

	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	room, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == errs.ErrNotFound {
			c.JSON(http.StatusNotFound, errorResp("not_found", "room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	c.JSON(http.StatusOK, roomToResponse(room))
}

// Deactivate godoc
// @Summary Деактивировать комнату
// @Description Помечает комнату как неактивную. Доступно только администратору
// @Tags admin-rooms
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комнаты"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/rooms/{id} [delete]
func (h *RoomsHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		return
	}

	if err := h.svc.Deactivate(c.Request.Context(), id); err != nil {
		if err == errs.ErrNotFound {
			c.JSON(http.StatusNotFound, errorResp("not_found", "room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	c.Status(http.StatusNoContent)
}


//маппер из domain в dto
func roomToResponse(r domain.Room) dto.RoomResponse {
	return dto.RoomResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Capacity:    r.Capacity,
		Location:    r.Location,
		IsActive:    r.IsActive,
	}
}