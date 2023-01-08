package routes

import (
	"rest-template/controller"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitDogRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	dogGroup := r.Group("/perro")
	{
		dogGroup.GET("/", controller.GetDog)
		//perroGroup.POST("/", postPerroHandler)
	}
}
