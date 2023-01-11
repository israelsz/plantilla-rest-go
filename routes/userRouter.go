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
		userGroup.POST("/register", controller.CreateUser)
		userGroup.GET("/:id", controller.GetUserByID)
		userGroup.GET("/email/:email", controller.GetUserByEmail)
		userGroup.GET("/", controller.GetAllUsers)
		userGroup.PUT("/update/:id", controller.UpdateUser)
		userGroup.DELETE("/delete/:id", controller.DeleteUser)
	}
}
