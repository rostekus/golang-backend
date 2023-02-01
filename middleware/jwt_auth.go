package middleware

import (
	"net/http"
	"rostekus/golang-backend/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		tokenString = strings.Split(tokenString, "Bearer ")[1]
		err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
