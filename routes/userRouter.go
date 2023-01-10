package routes

import (
	"rest-template/controller"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitUserRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	userGroup := r.Group("/user")
	{
		userGroup.POST("/create", controller.CreateUser)
		userGroup.GET("/get/:id", controller.GetUserByID)
		userGroup.GET("/email/:email", controller.GetUserByEmail)
		//perroGroup.POST("/", postPerroHandler)
	}
}
