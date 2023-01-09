package routes

import (
	"rest-template/controller"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitAuthRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	adminGroup := r.Group("/admin")
	adminGroup.Use(authMiddleware.MiddlewareFunc()) //Esta ruta solo puede ser utilizada por admins
	{
		adminGroup.GET("/", controller.GetAdmin)
		//perroGroup.POST("/", postPerroHandler)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authMiddleware.LoginHandler)
	}
}
