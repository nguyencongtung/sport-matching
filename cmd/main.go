package main

import (
	"app/internal/database"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/router"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to the database
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Serve static files from the "public" directory
	r.Static("/public", "./public")
	r.StaticFile("/style.css", "./public/style.css")
	r.StaticFile("/script.js", "./public/script.js")

	// Serve HTML files
	r.GET("/register.html", func(c *gin.Context) {
		c.File("./public/register.html")
	})
	r.GET("/login.html", func(c *gin.Context) {
		c.File("./public/login.html")
	})
	r.GET("/dashboard.html", func(c *gin.Context) {
		c.File("./public/dashboard.html")
	})

	// Redirect root to register page
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/register.html")
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userRepo)

	// Setup API routes
	api := r.Group("/api")
	router.AuthRoutes(api, authHandler)

	// Start the server
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
