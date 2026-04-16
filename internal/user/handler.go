package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	Registration(ctx context.Context, email string, password string) (*Response, error)
	Login(ctx context.Context, email string, password string) (*LoginResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var user Register
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}
	response, err := h.service.Registration(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
	}
	c.JSON(200, response)
}

func (h *Handler) Login(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}
	response, err := h.service.Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}
