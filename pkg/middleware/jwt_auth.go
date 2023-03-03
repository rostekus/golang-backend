package middleware

import (
	"net/http"
	"rostekus/golang-backend/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
		c.Abort()
		return
	}
	tokenString = strings.Split(tokenString, "Bearer ")[1]
	userID, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Set("user_id", userID)
	c.Next()
}
