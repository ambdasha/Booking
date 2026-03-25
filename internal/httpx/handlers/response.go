
//общий хэлпер для ответов

package handlers

import "github.com/gin-gonic/gin"

func errorResp(code, msg string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": msg}}
}