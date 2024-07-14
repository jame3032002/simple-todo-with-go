package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id     primitive.ObjectID `bson:"_id"`
	Task   string             `json:"task" binding:"required" bson:"task"`
	Status string             `json:"status" bson:"status"`
}

type UpdateTodo struct {
	Task   *string `json:"task,omitempty" bson:"task"`
	Status *string `json:"status,omitempty" validate:"oneof=todo doing done" bson:"status"`
}
