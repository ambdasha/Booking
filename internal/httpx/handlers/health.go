// быстрая проверка, что сервис жив и база доступна
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


type HealthHandler  struct{
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler{
	return &HealthHandler{db:db}
}

//проверка про бд невисит вечность
func (h *HealthHandler) Health(c *gin.Context){
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err:=h.db.Ping(ctx); err!=nil{
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ok":false, 
			"db": "down",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok":true})

}