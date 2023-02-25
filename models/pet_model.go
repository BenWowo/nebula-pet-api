package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PetSize struct {
	Height float64 `json:"height,omitempty"`
	Weight float64 `json:"weight,omitempty"`
}

type Pet struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"Name,omitempty" validate:"required"`
	DateOfBirth string             `json:"DateOfBirth,omitempty" validate:"required"`
	OwnerName   string             `json:"OwnerName,omitempty" validate:"required"`
	AnimalType  string             `json:"AnimalType,omitempty" validate:"required"`
	Breed       string             `json:"Breed,omitempty"`
	Size        PetSize            `json:"Size,omitempty"`
	FavoriteToy string             `json:"favorite_toy,omitempty"`
}
