package main

import (
	"log"

	"app/configs"
	"app/internal/database"
	"app/internal/handler"
	"app/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize configuration and logging
	configs.InitializeConfig()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})

	database.ConnectDB()
	app.Get("/ws", handler.ChatHandler, websocket.New(handler.WebSocketHandler))

	router.SetupRoutes(app)

	// Log server startup
	log.Fatal(app.Listen(":3000"))
}
