package main

import (
	"log"
	"net/http"
	"os"
	"rest-template/middleware"
	"rest-template/routes"
	"rest-template/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	// Se cargan variables de entorno
	utils.LoadEnv()

	// Log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Start template-go-rest")
	log.Printf("serverUp, %s ", os.Getenv("ADDR"))

	//Se fija el modo de gin desde las variables de entorno (debug | release)
	gin.SetMode(os.Getenv("GIN_MODE"))

	//Creacion de objeto gin
	app := gin.Default()
	// Cargar Cors
	app.Use(middleware.CorsMiddleware())

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	// Se registran las rutas(end-points) del proyecto
	routes.InitRoutes(app)

	//Se inicializa el servidor
	http.ListenAndServe(os.Getenv("ADDR"), app)

}
