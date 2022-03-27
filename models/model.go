package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoContent struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TodoTitle   string             `json:"TodoTitle,omitempty" bson:"TodoTitle,omitempty"`
	TodoContent string             `json:"TodoContent,omitempty" bson:"TodoContent,omitempty"`
}
