package router

import (
	"net/http"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/image"
	"rostekus/golang-backend/pkg/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewImageServiceRouter(imageHandler *image.Handler, healthHandler *health.Handler) *Router {
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
	public.POST("/upload", imageHandler.UploadFile)
	public.POST("/download", imageHandler.DownloadFile)

	private := router.Group("api/v1")
	private.GET("/checkjwt", middleware.JWTAuth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the private route!"})
	})

	return &Router{Router: router}
}
