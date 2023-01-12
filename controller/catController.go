package controller

import (
	"log"
	"net/http"
	"rest-template/models"
	"rest-template/services"

	"github.com/gin-gonic/gin"
)

// Servicio que permite Crear un gato
func CreateCat(ctx *gin.Context) {
	// Obtiene los datos del gato a partir del cuerpo de la solicitud HTTP.
	var cat models.Cat
	//Si ocurrio un error
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Se llama al servicio que crea el gato en la base de datos Mongo
	createdCat, err := services.CreateCatService(cat)
	// Si ocurrio un error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al crear el gato"})
		return
	}
	log.Println("Gato creado en la base de datos")
	ctx.JSON(http.StatusCreated, createdCat)

}

// Servicio que permite retornar un gato buscandolo por su id
func GetCatByID(ctx *gin.Context) {
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	catID := ctx.Param("id")
	// Se busca al gato en la base de datos
	resultCat, err := services.GetCatByIDService(catID)
	// Si ocurrio un error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener al gato"})
		return
	}
	// Devuelve el gato encontrado.
	ctx.JSON(http.StatusOK, resultCat)
}

// Servicio que permite eliminar a un gato de la base de datos
func DeleteCat(ctx *gin.Context) {
	catID := ctx.Param("id")
	err := services.DeleteCatService(catID)
	// Si hubo un error
	if err != nil {
		// Ocurrió un error durante la búsqueda.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Devuelve el mensaje de confirmación
	ctx.JSON(http.StatusOK, "Gato eliminado")
}

// Servicio que permite actualizar a un gato
func UpdateCat(ctx *gin.Context) {
	//Se crea un modelo de gato
	var updatedCat models.Cat
	//Se guardan los datos de la petición http en el modelo
	//Si hubo un error
	if err := ctx.ShouldBindJSON(&updatedCat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos del Gato"})
		return
	}
	//Se consigue el id del gato a actualizar
	catID := ctx.Param("id")
	//Se actualiza el gato en la base de datos
	updatedCat, err := services.UpdateCatService(updatedCat, catID)
	//Si ocurrio un error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al actualizar al gato"})
		return
	}
	//Se envia la respuesta http
	ctx.JSON(http.StatusCreated, updatedCat)
}

// Servicio para conseguir a todos los gatos de la base de datos
func GetAllCats(ctx *gin.Context) {
	// Se consiguen a los gatos de la base de datos
	resultCats, err := services.GetAllCatsService()
	//Si ocurrio un error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener los gatos"})
		return
	}
	//Se envia la respuesta http
	ctx.JSON(http.StatusCreated, resultCats)
}
