package main

import (
	"app/configs"
	"app/internal/database"
	"app/internal/router"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize global logger
	configs.LoggerInit("auth_server3.log", slog.LevelDebug)

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})

	//database.ConnectDB()
	database.ConnectMongoDB()

	router.SetupRoutes(app)

	// Log server startup
	slog.Info(app.Listen(":3000").Error())

}
