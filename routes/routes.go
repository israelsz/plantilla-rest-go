package routes

import (
	"github.com/gin-gonic/gin"
)

// InitRoutes registra las rutas junto a las funciones que ejecutan
func InitRoutes(r *gin.Engine) {
	// Registra las rutas del grupo de perros del archivo perroRouter.go
	InitDogRoutes(r)
	//Registra las rutas del grupo de gatos del archivo gato.go en package Routes
	InitCatRoutes(r)
	//Registra las rutas del grupo de gatos del archivo gatoRouter.go
	InitAuthRoutes(r)
	//Registra las rutas del grupo de usuario del archivo usuarioRouter.go
	InitUserRoutes(r)
}
