package routes

import (
	"rest-template/controller"
	"rest-template/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitCatRoutes(r *gin.Engine) {
	// Define a group of routes with a shared set of middleware
	// Se define un grupo de rutas
	catGroup := r.Group("/cat")
	{
		//Solo un usuario logueado sin importar su rol, puede crear un gato
		catGroup.POST("/", middleware.LoadJWTAuth().MiddlewareFunc(), controller.CreateCat)
		catGroup.GET("/:id", controller.GetCatByID)
		catGroup.GET("/", controller.GetAllCats)
		//Solo un usuario o admin logueados pueden actualizar a un gato
		catGroup.PUT("/:id", middleware.SetRoles(middleware.RolAdmin, middleware.RolUser), middleware.LoadJWTAuth().MiddlewareFunc(), controller.UpdateCat)
		//Solo un Admin logueado puede eliminar a un gato
		catGroup.DELETE("/:id", middleware.SetRoles(middleware.RolAdmin), middleware.LoadJWTAuth().MiddlewareFunc(), controller.DeleteCat)
	}
}
