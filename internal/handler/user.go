package handler

import (
	"app/internal/database"
	"app/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUser retrieves a user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := database.DB
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserProfile updates a user's profile information
func UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	db := database.DB
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind only the profile-related fields
	var profileUpdate struct {
		Names             string `json:"names"`
		Gender            string `json:"gender"`
		DateOfBirth       string `json:"date_of_birth"`
		Bio               string `json:"bio"`
		Interests         string `json:"interests"`
		LookingFor        string `json:"looking_for"`
		ProfilePictureURLs model.StringArray `json:"profile_picture_urls"` // Changed to StringArray
		Location          string `json:"location"`
		DistancePreference int    `json:"distance_preference"`
	}

	if err := c.ShouldBindJSON(&profileUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Names = profileUpdate.Names
	user.Gender = profileUpdate.Gender
	user.DateOfBirth = profileUpdate.DateOfBirth
	user.Bio = profileUpdate.Bio
	user.Interests = profileUpdate.Interests
	user.LookingFor = profileUpdate.LookingFor
	user.ProfilePictureURLs = profileUpdate.ProfilePictureURLs // Changed to ProfilePictureURLs
	user.Location = profileUpdate.Location
	user.DistancePreference = profileUpdate.DistancePreference

	db.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully", "user": user})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := database.DB
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
