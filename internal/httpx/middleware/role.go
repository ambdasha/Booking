package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


//защита эндпоинтов по роли
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context){
		_, currentRole := MustUser(c)
		if currentRole !=role{
			c.AbortWithStatusJSON(http.StatusForbidden,gin.H{
				"error":gin.H{"code": "forbidden", "message": "insufficient permissions"},
			})
			return
		}
		c.Next()
	}

}