package routes

import (
	"rest-template/controller"
	"rest-template/middleware"

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
		// Solo Usuarios y Admins logueados pueden actualizar datos de usuario
		userGroup.PUT("/:id", middleware.SetRoles(middleware.RolAdmin, middleware.RolUser), middleware.LoadJWTAuth().MiddlewareFunc(), controller.UpdateUser)
		// Solo Admins logueados pueden borrar usuarios
		userGroup.DELETE("/:id", middleware.SetRoles(middleware.RolAdmin), middleware.LoadJWTAuth().MiddlewareFunc(), controller.DeleteUser)
	}
}
