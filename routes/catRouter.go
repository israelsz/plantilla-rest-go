package routes

import (
	"rest-template/controller"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitCatRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	catGroup := r.Group("/cat")
	{
		catGroup.POST("/", controller.CreateCat)
		catGroup.GET("/:id", controller.GetCatByID)
		catGroup.GET("/", controller.GetAllCats)
		catGroup.DELETE("/:id", controller.DeleteCat)
		catGroup.PUT("/:id", controller.UpdateCat)
	}
}
