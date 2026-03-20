package middleware

import(
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//вешает на каждый HTTP-запрос уникальный request id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-Id")

		if id == "" {
			id = uuid.NewString()
		}

		c.Writer.Header().Set("X-Request-Id", id)
		c.Set("request_id", id)
		c.Next()
	}
}