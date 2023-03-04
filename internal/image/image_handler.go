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
	userID := c.GetString("user_id")
	res, err := h.service.UploadFileToMongoDB(c, &req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.InsertImageDataToDB(c, &req, userID, res.ImageID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	imageMessage := ImageMessage{
		ImageID: res.ImageID,
		UserID:  userID,
	}
	err = h.service.PublishMessage(imageMessage)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DownloadFile(c *gin.Context) {
	var req ImageDownloadRequest
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("user_id").(string)
	resp, fileBuffer, err := h.service.DownloadFile(c, &req, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := c.Writer.Write(fileBuffer.Bytes()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", http.DetectContentType(fileBuffer.Bytes()))

	c.Header("Content-Length", strconv.Itoa(fileBuffer.Len()))
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetImagesForUser(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	images, err := h.service.GetImagesForUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
