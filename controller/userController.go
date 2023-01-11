package controller

import (
	"log"
	"net/http"
	"rest-template/models"
	"rest-template/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Se establecen los nombres de la colección que se traeran desde la base de datos
*/
const (
	CollectionNameUser = "User"
)

// Función para crear un usuario e insertarlo a la base de datos de mongodb

func CreateUser(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos de usuario"})
		return
	}

	createdUser, err := services.CreateUserService(newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

// Función para obtener un usuario por id
func GetUserByID(ctx *gin.Context) {
	// Obtiene el ID del usuario a partir del parámetro de la ruta.
	userID := ctx.Param("id")
	resultUser, err := services.GetUserByIDService(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

// Función para obtener un usuario por id
func GetUserByEmail(ctx *gin.Context) {
	userEmail := ctx.Param("email")
	resultUser, err := services.GetUserByEmailService(userEmail)
	// Si hubo un error
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			log.Println("Usuario no encontrado")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cat not found"})
			return
		}
		// Ocurrió un error durante la búsqueda.
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Si no hay error se retorna el usuario
	log.Println("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	ctx.JSON(http.StatusOK, resultUser)
}

func GetAllUsers(ctx *gin.Context) {
	resultUser, err := services.GetAllUserService()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, resultUser)
}

func UpdateUser(ctx *gin.Context) {
	var updatedUser models.User
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos de usuario"})
		return
	}
	userId := ctx.Param("id")
	updatedUser, err := services.UpdateUserService(updatedUser, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al updatear al usuario"})
		return
	}
	ctx.JSON(http.StatusCreated, updatedUser)
}

// Función para obtener un usuario por id
func DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := services.DeleteUserService(userID)
	// Si hubo un error
	if err != nil {
		// Ocurrió un error durante la búsqueda.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Si no hay error se retorna el usuario
	log.Println("Se elimino el usuario")
	// Devuelve el usuario encontrado.
	ctx.JSON(http.StatusOK, "Usuario eliminado")
}
