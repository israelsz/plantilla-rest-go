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
		userGroup.POST("/", controller.CreateUser)
		userGroup.GET("/:id", controller.GetUserByID)
		userGroup.GET("/email/:email", controller.GetUserByEmail)
		userGroup.GET("/", controller.GetAllUsers)
		userGroup.PUT("/:id", controller.UpdateUser)
		userGroup.DELETE("/:id", controller.DeleteUser)
	}
}
