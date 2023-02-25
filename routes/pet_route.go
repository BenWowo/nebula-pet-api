package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func PetRoute(router *gin.Engine) {
	// All routes related to pets comes here
	router.POST("/pet", controllers.CreatePet())
	router.GET("/pet/:petId", controllers.GetAPet())
	router.PUT("/pet/:petId", controllers.EditAPet())
	router.DELETE("/pet/:petId", controllers.DeleteAPet())
	router.GET("/pets", controllers.GetAllPets())
}
