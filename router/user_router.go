package router

import (
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var userRouter *gin.Engine

func NewUserServiceRouter(userHandler *user.Handler, healthHandler *health.Handler) *gin.Engine {
	userRouter = gin.Default()

	userRouter.Use(cors.New(cors.Config{
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
	public := userRouter.Group("/")

	public.POST("/signup", userHandler.CreateUser)
	public.POST("/login", userHandler.LoginUser)
	public.GET("/health", healthHandler.HealthCheck)

	return userRouter
}
