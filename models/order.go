package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	RecetteRefer int                `json:"Recette_id"`
	Recette      Recette            `gorm:"foreignKey:RecetteRefer"`
	UserRefer    int                `json:"user_id"`
	User         User               `gorm:"foreignKey:UserRefer"`
}
