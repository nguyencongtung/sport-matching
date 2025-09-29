package handler

import (
	"app/internal/database"
	"app/internal/models"
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var messageCollection *mongo.Collection = database.GetCollection(database.ConnectMongoDB(), "messages")

// Hello handle api status
func Hello(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("API is healthy")

	newMsg := models.Message{
		Id:       primitive.NewObjectID(),
		Name:     "John Doe",
		Location: "New York",
		Title:    "Hello World",
	}

	// Uncomment and use the collection variable if needed, or remove the line if it's unnecessary.
	collection, err := messageCollection.InsertOne(ctx, newMsg)
	if err != nil {
		slog.Error("Error inserting document:", slog.Any("error", err))
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't insert document", "data": err})
	}
	slog.Info("Inserted document ID:", slog.Any("id", collection.InsertedID))
	_ = collection
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok! ", "data": collection})
}
