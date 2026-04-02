//регистрация и логин на уровне HTTP слоя

package handlers

import (
	"context"
	"net/http"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (domain.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error)
}

type AuthHandler struct {
	svc       AuthService
	validator *validator.Validate
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{
		svc:       svc,
		validator: validator.New(),
	}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	u, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if err == errs.ErrConflict {
			c.JSON(http.StatusConflict, errorResp("email_taken", "email already registered"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResp("internal_error", "something went wrong"))
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	})
}

// Login godoc
// @Summary Вход в систему
// @Description Возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid json"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("validation_error", "invalid fields"))
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResp("invalid_credentials", "wrong email or password"))
		return
	}

	c.JSON(http.StatusOK, token)
}