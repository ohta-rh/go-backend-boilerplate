package handler

import (
	"errors"
	"net/http"
	"strconv"

	domainUser "easy-go-backend/internal/domain/user"
	"easy-go-backend/internal/schema"
	"easy-go-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	userInteractor *usecase.UserInteractor
}

// NewUserHandler creates a new user handler.
func NewUserHandler(ui *usecase.UserInteractor) *UserHandler {
	return &UserHandler{
		userInteractor: ui,
	}
}

// GetUser handles GET request for a single user.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userInteractor.GetUser(c, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, schema.FromEntity(user))
}

// CreateUser handles POST request to create a user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req schema.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEntity := req.ToEntity()

	createdUser, err := h.userInteractor.CreateUser(c, userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schema.FromEntity(createdUser))
}

// UpdateUser handles PUT request to update a user.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req schema.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEntity := req.ToEntity()
	userEntity.ID = id

	updatedUser, err := h.userInteractor.UpdateUser(c, userEntity)
	if err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, schema.FromEntity(updatedUser))
}

// DeleteUser handles DELETE request to delete a user.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userInteractor.DeleteUser(c, id); err != nil {
		if errors.Is(err, domainUser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsers handles GET request to retrieve all users.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userInteractor.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema.FromEntities(users))
}

// RegisterRoutes はユーザー関連のルートを登録します.
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.GET("", h.GetAllUsers)
	users.POST("", h.CreateUser)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}
