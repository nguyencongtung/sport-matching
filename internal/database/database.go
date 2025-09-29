package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

// Client instance
var mongoDB *mongo.Client
