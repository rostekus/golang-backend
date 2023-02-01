package router

import (
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/user"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func NewRouter(userHandler *user.Handler, healthHandler *health.Handler) *gin.Engine {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.LoginUser)
	r.GET("health", healthHandler.HealthCheck)
	return r
}
