package routes

import (
	"rest-template/controller"
	"rest-template/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitAuthRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	adminGroup := r.Group("/loggeados")
	{
		adminGroup.GET("/admin", middleware.SetRoles(middleware.RolAdmin), middleware.LoadJWTAuth().MiddlewareFunc(), controller.GetAdmin)                       // Esta ruta solo puede ser usada por admins
		adminGroup.GET("/usuario", middleware.SetRoles(middleware.RolAdmin, middleware.RolUser), middleware.LoadJWTAuth().MiddlewareFunc(), controller.GetAdmin) // Esta ruta solo puede ser usada por admins y usuarios loggeados
		//perroGroup.POST("/", postPerroHandler)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", middleware.LoadJWTAuth().LoginHandler)
	}
}
