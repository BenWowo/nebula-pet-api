package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var petCollection *mongo.Collection = configs.GetCollection(configs.DB, "pets")
var validate = validator.New()

func CreatePet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var pet models.Pet
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusBadRequest, responses.PetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newPet := models.Pet{
			Id:          primitive.NewObjectID(),
			Name:        pet.Name,
			DateOfBirth: pet.DateOfBirth,
			OwnerName:   pet.OwnerName,
			AnimalType:  pet.AnimalType,
			Breed:       pet.Breed,
			Size: models.PetSize{
				Height: pet.Size.Height,
				Weight: pet.Size.Weight,
			},
			FavoriteToy: pet.FavoriteToy,
		}

		result, err := petCollection.InsertOne(ctx, newPet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.PetResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		petId := c.Param("petId")
		var pet models.Pet
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(petId)

		err := petCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&pet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.PetResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": pet}})
	}
}

func EditAPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		petId := c.Param("petId")
		var pet models.Pet
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(petId)

		//validate the request body
		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusBadRequest, responses.PetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		//update := bson.M{"name": pet.Name, "DateOfBirth": pet.DateOfBirth, "OwnerName": pet.OwnerName, "AnimalType": pet.AnimalType, "Breed": pet.Breed, "Height": pet.Height, "Weight": pet.Weight, "FavoriteToy": pet.FavoriteToy}
		update := bson.M{"name": pet.Name, "DateOfBirth": pet.DateOfBirth, "OwnerName": pet.OwnerName, "AnimalType": pet.AnimalType, "Breed": pet.Breed, "Height": pet.Size.Height, "Weight": pet.Size.Weight, "FavoriteToy": pet.FavoriteToy}

		//
		result, err := petCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated pet details
		var updatedPet models.Pet
		if result.MatchedCount == 1 {
			err := petCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPet)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.PetResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedPet}})
	}
}

func DeleteAPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		petId := c.Param("petId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(petId)

		result, err := petCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.PetResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Pet with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.PetResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Pet successfully deleted!"}},
		)
	}
}

func GetAllPets() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var pets []models.Pet
		defer cancel()

		results, err := petCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singlePet models.Pet
			if err = results.Decode(&singlePet); err != nil {
				c.JSON(http.StatusInternalServerError, responses.PetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			pets = append(pets, singlePet)
		}

		c.JSON(http.StatusOK,
			responses.PetResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": pets}},
		)
	}
}
