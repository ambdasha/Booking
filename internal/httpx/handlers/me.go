
//проверить, что JWT middleware работает и узнать свой id и role без отдельного запроса в БД

package handlers

import (
	"net/http"

	"booking/internal/httpx/middleware"

	"github.com/gin-gonic/gin"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler { return &MeHandler{} }

// Me godoc
// @Summary Текущий пользователь
// @Description Возвращает данные текущего авторизованного пользователя
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /me [get]
func (h *MeHandler) Me(c *gin.Context) {
	uid, role := middleware.MustUser(c)
	c.JSON(http.StatusOK, gin.H{
		"id":   uid,
		"role": role,
	})
}