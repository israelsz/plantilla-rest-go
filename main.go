package main

import (
	"log"
	"net/http"
	"os"
	"rest-template/config"
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

	// Conectar la base de datos
	config.LoadDatabase()
	println("Conexion lograda")

	//Creacion de objeto gin
	app := gin.Default()
	// Cargar Cors
	app.Use(middleware.CorsMiddleware())

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	// Se inicia el middleware JWT
	//authMiddleware := middleware.LoadJWTAuth()

	// Se registran las rutas(end-points) del proyecto
	routes.InitRoutes(app)

	//Se inicializa el servidor
	http.ListenAndServe(os.Getenv("ADDR"), app)

}
