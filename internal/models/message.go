package models

// type Message struct {
// 	Room    string `json:"room"`
// 	Sender  string `json:"sender"`
// 	Content string `json:"content"`
// }

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Title    string             `json:"title,omitempty" validate:"required"`
}
