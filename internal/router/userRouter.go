package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up user-related routes
func UserRoutes(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware()) // Apply authentication middleware to user routes
	{
		user.GET("/:id", handler.GetUser)
		user.PUT("/profile/:id", handler.UpdateUserProfile)
		user.DELETE("/:id", handler.DeleteUser)
	}
}
