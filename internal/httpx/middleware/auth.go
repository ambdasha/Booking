package middleware

import (
	"net/http"
	"strings"

	"booking/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserIDKey = "user_id"
	ctxRoleKey   = "role"
)

//JWT-middleware
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "missing token"},
			})
			return
		}

		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.ParseJWT(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "unauthorized", "message": "invalid token"},
			})
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxRoleKey, claims.Role)

		c.Next()
	}
}

//хелпер, который достаёт user_id и role из gin.Context
func MustUser(c *gin.Context) (int64, string) {
	uidAny, _ := c.Get(ctxUserIDKey)
	roleAny, _ := c.Get(ctxRoleKey)

	uid, _ := uidAny.(int64)
	role, _ := roleAny.(string)

	return uid, role
}