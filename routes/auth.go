package routes

import (
	"rest-template/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", middleware.LoadJWTAuth().LoginHandler)
		authGroup.POST("/refresh_token", middleware.LoadJWTAuth().RefreshHandler)
		authGroup.POST("/logout", middleware.LoadJWTAuth().LogoutHandler)

	}
}
