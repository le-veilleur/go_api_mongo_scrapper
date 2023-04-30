package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recette struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Descriptions string             `json:"descriptions" bson:"descriptions"`
	Ingredients  string             `json:"ingredients" bson:"ingredients"`
	Photos       string             `json:"photos" bson:"photos"`
	Instructions   string             `json:"instructions" bson:"instructions"`
	Page         string             `json:"line" bson:"line"`
	SerialNumber string             `json:"serial_number" bson:"serial_number"`
}
