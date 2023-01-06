package routes

import (
	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitCatRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	catGroup := r.Group("/gato")
	{
		//gatoGroup.GET("/", getGatoHandler)
		//gatoGroup.POST("/", postGatoHandler)
	}
}
