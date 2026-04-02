package handlers
//HTTP-слой для блокировок комнат доспупных админу
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

//связь между handler и service
type BlockService interface {
	Create(ctx context.Context, adminID int64, roomID int64, req dto.CreateBlockRequest) (domain.RoomBlock, error)
	Delete(ctx context.Context, blockID int64) error
}

type BlocksHandler struct {
	svc       BlockService
	validator *validator.Validate
}

func NewBlocksHandler(svc BlockService) *BlocksHandler {
	return &BlocksHandler{svc: svc, validator: validator.New()}
}

// Create godoc
// @Summary Создать блокировку комнаты
// @Description Создает блокировку комнаты на указанный интервал. Доступно только администратору
// @Tags admin-blocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комнаты"
// @Param input body dto.CreateBlockRequest true "Данные блокировки"
// @Success 201 {object} dto.BlockResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /admin/rooms/{id}/blocks [post]
func (h *BlocksHandler) Create(c *gin.Context) {
	adminID, _ := middleware.MustUser(c)

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64) //парсинг

	if err != nil || roomID <= 0 {  // если id<=0 то ошибка
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid room id"))
		return
	}

	var req dto.CreateBlockRequest

	if err := c.ShouldBindJSON(&req); err != nil { //если json кривой, то ошибка
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	if err := h.validator.Struct(req); err != nil {//проверка тегов validate:"..." в DTO
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}
// вызов метода из service
	b, err := h.svc.Create(c.Request.Context(), adminID, roomID, req)

	if err != nil {
		switch err { //обработка ошибок
		case errs.ErrValidation:
			c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid interval"))
		case errs.ErrForbidden:
			c.JSON(http.StatusForbidden, errorResp("forbidden", "room is inactive"))
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, errorResp("not_found", "room not found"))
		case errs.ErrConflict:
			c.JSON(http.StatusConflict, errorResp("block_conflict", "block overlaps with existing block"))
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		}
		return
	}
//возвращается dto
	c.JSON(http.StatusCreated, dto.BlockResponse{
		ID:        b.ID,
		RoomID:    b.RoomID,
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Reason:    b.Reason,
		CreatedAt: b.CreatedAt,
	})
}


// Delete godoc
// @Summary Удалить блокировку
// @Description Удаляет блокировку комнаты. Доступно только администратору
// @Tags admin-blocks
// @Produce json
// @Security BearerAuth
// @Param block_id path int true "ID блокировки"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/blocks/{block_id} [delete]
func (h *BlocksHandler) Delete(c *gin.Context) {

	blockID, err := strconv.ParseInt(c.Param("block_id"), 10, 64)
	if err != nil || blockID <= 0 { // id<=0 -> err
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid block id"))
		return
	}

// вызов метода из service
	if err := h.svc.Delete(c.Request.Context(), blockID); err != nil {

		switch err {
		case errs.ErrValidation:
			c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid id"))
		case errs.ErrNotFound:
			c.JSON(http.StatusNotFound, errorResp("not_found", "block not found"))
		default:
			c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		}
		return
	}

	c.Status(http.StatusNoContent)// 204 No Content
}