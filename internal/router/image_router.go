package router

import (
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
	router.POST("/images", middleware.JWTAuth, imageHandler.UploadFile)
	router.GET("/images/:id", middleware.JWTAuth, imageHandler.DownloadFile)
	router.GET("/images", middleware.JWTAuth, imageHandler.GetImagesForUser)
	return &Router{Router: router}
}
