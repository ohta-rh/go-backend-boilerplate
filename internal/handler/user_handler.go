package handler

import (
	"net/http"
	"strconv"

	"easy-go-backend/internal/domain/user"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST request to create a user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userDTO user.User
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userInteractor.CreateUser(c, &userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser handles PUT request to update a user.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userDTO user.User
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userDTO.ID = id

	updatedUser, err := h.userInteractor.UpdateUser(c, &userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE request to delete a user.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userInteractor.DeleteUser(c, id); err != nil {
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

	c.JSON(http.StatusOK, users)
}
