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
		catGroup.POST("/create", controller.CreateCat)
		catGroup.GET("find/:id", controller.GetCatByID)
		catGroup.GET("getAll", controller.GetAllCats)
		catGroup.DELETE("delete/:id", controller.DeleteCat)
		catGroup.PUT("update/:id", controller.UpdateCat)
	}
}
