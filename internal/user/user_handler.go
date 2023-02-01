package user

import (
	"net/http"

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

func (h *Handler) CreateUser(c *gin.Context) {
	var u CreateUserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) LoginUser(c *gin.Context) {
	var u LoginUserRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	res, err := h.service.LoginUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, res)
}
