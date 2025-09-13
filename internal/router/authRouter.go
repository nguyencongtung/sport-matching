package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up authentication routes
func AuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)

		// Protected routes
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", authHandler.GetCurrentUser)
			protected.PUT("/profile", userHandler.UpdateUserProfile) // New route for profile setup
		}
	}
}
