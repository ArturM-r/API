package APIs

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, title string) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*API, error)
	GetAll(ctx context.Context, userID uuid.UUID, filter QueryFilterFromHandler) ([]API, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, title *string, completed *bool) error
}

type APHandler struct {
	service Service
}

func NewAPIHandler(service Service) *APHandler {
	return &APHandler{service: service}
}

func (h *APHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := value.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}
	var req Create

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind create API request body"})
		return
	}
	if err := h.service.Create(ctx, userID, req.Title); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "API created successfully",
	})
}

func (h *APHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := value.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(ctx, id, userID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "API deleted successfully",
	})
}
func (h *APHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	api, err := h.service.GetByID(ctx, id, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api)
}

func (h *APHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	var filter QueryFilterFromHandler
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := value.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}
	apis, err := h.service.GetAll(ctx, userID, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apis)
}
func (h *APHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req Update
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind update API request body"})
		return
	}

	if err := h.service.Update(ctx, id, userID, req.Title, req.Completed); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API updated successfully",
	})
}
