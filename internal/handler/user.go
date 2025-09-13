package handler

import (
	"app/internal/database"
	"app/internal/model"
	"app/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler struct holds dependencies for user-related handlers
type UserHandler struct {
	UserRepository repository.UserRepository
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: userRepo}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserRepository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// CreateUser creates a new user (this functionality is now handled by AuthHandler.Register)
// This function is kept as a placeholder or for other non-auth related user creation if needed.
func (h *UserHandler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "User creation is handled by /api/auth/register"})
}

// UpdateUser updates a user's information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserRepository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UserRepository.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// ProfileSetupRequest defines the structure for profile setup data
type ProfileSetupRequest struct {
	FirstName        string   `json:"firstName"`
	DateOfBirth      string   `json:"dateOfBirth"`
	Gender           string   `json:"gender"`
	ShowGender       bool     `json:"showGender"`
	DistancePreference int      `json:"distancePreference"`
	ProfilePictures  []string `json:"profilePictures"`
}

// UpdateUserProfile handles updating user profile information during the setup process
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var req ProfileSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming user ID is available from JWT token or session
	// For now, let's use a placeholder user ID or retrieve from context if authenticated
	// This part needs proper authentication middleware to extract user ID
	userID := c.MustGet("userID").(uint) // Example: Get userID from authenticated context

	user, err := h.UserRepository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.FirstName = req.FirstName
	user.DateOfBirth = req.DateOfBirth
	user.Gender = req.Gender
	user.ShowGender = req.ShowGender
	user.DistancePreference = req.DistancePreference

	// Convert []string to JSON string for ProfilePictures
	if len(req.ProfilePictures) > 0 {
		picturesJSON, err := json.Marshal(req.ProfilePictures)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal profile pictures"})
			return
		}
		user.ProfilePictures = string(picturesJSON)
	}

	if err := h.UserRepository.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully", "user": user})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := database.DB // Direct DB access, consider using repository for consistency
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
