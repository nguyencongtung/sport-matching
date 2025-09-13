package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up authentication routes
func AuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)

		// Protected routes
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
	}
}
