package image

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service
}

func NewHandler(s *service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) UploadFile(c *gin.Context) {
	var req ImageUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.UploadFileToMongoDB(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DownloadFile(c *gin.Context) {
	var req ImageDownloadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileBuffer, err := h.service.DownloadFile(&req.ImageName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := c.Writer.Write(fileBuffer.Bytes()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", http.DetectContentType(fileBuffer.Bytes()))

	// Set the Content-Length header based on the length of the buffer
	c.Header("Content-Length", strconv.Itoa(fileBuffer.Len()))
	c.JSON(http.StatusOK,
		gin.H{})
}
