package handler

import (
	config "app/configs"
	"app/internal/model"
	"app/internal/repository"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(config.Config("JWT_SECRET"))

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo repository.UserRepository
	validate *validator.Validate
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		validate: validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	if err := h.validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user with email already exists
	existingUser, err := h.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := h.userRepo.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate login request
	if err := h.validate.Struct(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	SetAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// GetCurrentUser retrieves the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// This handler will be protected by JWT middleware
	// The user ID will be extracted from the JWT token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.GetUserByID(userID.(uint))
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
		}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout handles user logout (optional, as JWTs are stateless)
func (h *AuthHandler) Logout(c *gin.Context) {
	ClearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken handles refreshing JWT tokens (optional)
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Refresh token not implemented"})
}

// ForgotPassword handles password reset requests (optional)
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Forgot password not implemented"})
}

// ResetPassword handles password reset confirmation (optional)
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Reset password not implemented"})
}

// VerifyEmail handles email verification (optional)
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Email verification not implemented"})
}

// UpdateProfile handles updating user profile information (optional)
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update profile not implemented"})
}

// DeleteAccount handles deleting user accounts (optional)
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete account not implemented"})
}

// ChangePassword handles changing user password (optional)
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Change password not implemented"})
}

// GetUsers retrieves a list of all users (admin only, optional)
func (h *AuthHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Get users not implemented"})
}

// GetUser retrieves a single user by ID (admin only, optional)
func (h *AuthHandler) GetUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Get user not implemented"})
}

// UpdateUser updates a user's information (admin only, optional)
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update user not implemented"})
}

// DeleteUser deletes a user (admin only, optional)
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete user not implemented"})
}

// GenerateToken generates a JWT token for a given user
func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // Token valid for 7 days
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["userID"].(float64))
		return userID, nil
	}
	return 0, fmt.Errorf("invalid token")
}

// SetAuthCookie sets an authentication cookie
func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie("token", token, int(time.Hour*24*7/time.Second), "/", "localhost", false, true)
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
}
