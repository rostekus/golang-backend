package router

import (
	"fmt"
	"net/http"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/user"
	"rostekus/golang-backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewUserServiceRouter(userHandler *user.Handler, healthHandler *health.Handler) *Router {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
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
	public := router.Group("/")
	public.POST("/signup", userHandler.CreateUser)
	public.POST("/login", userHandler.LoginUser)
	public.GET("/health", healthHandler.HealthCheck)

	private := router.Group("api/v1")
	private.GET("/checkjwt", middleware.JWTAuth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome to the private route! id: %v", c.MustGet("user_id"))})
	})

	return &Router{Router: router}
}
