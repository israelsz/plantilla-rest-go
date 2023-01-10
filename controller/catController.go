package controller

import (
	"context"
	"log"
	"net/http"
	"rest-template/config"
	"rest-template/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionNameCat = "Cat"
)

type CatController struct {
	DbConnection *config.DbConnection
}

// NewCatController crea una nueva instancia de CatController.
func NewCatController() *CatController {
	// Establece una conexión a la base de datos de MongoDB.
	dbConnection := config.NewDbConnection()
	log.Println("Conexion establecida")

	return &CatController{
		DbConnection: dbConnection,
	}
}

func CreateCat(ctx *gin.Context) {
	log.Println("Service: CreateCat")
	// Crea una nueva instancia de CatController.
	c := NewCatController()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer c.DbConnection.Close()

	// Obtiene la colección de gatos.
	collection := c.DbConnection.GetCollection(CollectionNameCat)

	// Obtiene los datos del gato a partir del cuerpo de la solicitud HTTP.
	var cat models.Cat
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Genera un nuevo ID único para el gato.
	cat.ID = primitive.NewObjectID()
	// Establece la fecha de creación y actualización del gato.
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	// Inserta el gato en la colección.
	_, err := collection.InsertOne(context.TODO(), cat)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Devuelve el gato creado.
	ctx.JSON(http.StatusCreated, cat)
	log.Println("Gato creado")
}

func GetCatByID(ctx *gin.Context) {
	log.Println("Service: GetCatByID")
	// Crea una nueva instancia de CatController.
	c := NewCatController()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer c.DbConnection.Close()
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	catID := ctx.Param("id")

	// Crea un objeto ID de MongoDB a partir del ID del gato.
	oid, err := primitive.ObjectIDFromHex(catID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"_id": oid}

	// Obtiene la colección de gatos.
	collection := c.DbConnection.GetCollection(CollectionNameCat)

	// Ejecuta el método FindOne() con el filtro y el contexto.
	var result models.Cat
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cat not found"})
			return
		}
		// Ocurrió un error durante la búsqueda.
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devuelve el gato encontrado.
	ctx.JSON(http.StatusOK, result)
}
func DeleteCat(ctx *gin.Context) {
	log.Println("Service: DeleteCat")
	// Crea una nueva instancia de CatController.
	c := NewCatController()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer c.DbConnection.Close()

	// Obtiene la colección de gatos.
	collection := c.DbConnection.GetCollection(CollectionNameCat)

	// Obtiene el ID del gato a eliminar a partir de la URL de la solicitud HTTP.
	catID := ctx.Param("id")

	// Convierte el ID del gato a un ObjectID de MongoDB.
	objID, err := primitive.ObjectIDFromHex(catID)
	if err != nil {
		// Si ocurre un error al convertir el ID, devuelve un error HTTP 500 con el mensaje de error.
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Crea el filtro para eliminar el gato con el ID especificado.
	filter := bson.M{"_id": objID}

	// Elimina el gato de la colección.
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// Si ocurre un error al eliminar el gato, devuelve un error HTTP 500 con el mensaje de error.
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("Gato eliminado")
	// Devuelve una respuesta vacía con código HTTP 200 para indicar que la eliminación se ha realizado correctamente.
	ctx.Status(http.StatusOK)
}

func UpdateCat(ctx *gin.Context) {
	log.Println("Service: UpdateCat")
	// Crea una nueva instancia de CatController.
	// Crea una nueva instancia de CatController.
	c := NewCatController()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer c.DbConnection.Close()

	// Obtiene la colección de gatos.
	collection := c.DbConnection.GetCollection(CollectionNameCat)

	// Obtiene el ID del gato a partir de la URL de la solicitud HTTP.
	catID := ctx.Param("id")
	// Convierte la cadena hexadecimal en un ObjectID.
	objectID, err := primitive.ObjectIDFromHex(catID)
	if err != nil {
		// Ha ocurrido un error al convertir la cadena hexadecimal en un ObjectID
	}

	// Obtiene los datos del gato a partir del cuerpo de la solicitud HTTP.
	var cat models.Cat
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		// Ha ocurrido un error al obtener los datos del gato
	}

	// Establece la fecha de actualización del gato.
	cat.UpdatedAt = time.Now()

	// Actualiza el gato en la colección.
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": cat}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// Ha ocurrido un error al actualizar el gato
	}

	// Devuelve el gato actualizado.
	ctx.Status(http.StatusOK)
	log.Println("Gato actualizado")
}

// Función para traer a todos los gatos de la base de datos
func GetAllCats(ctx *gin.Context) {
	log.Println("Service: GetAllCats")
	// Crea una nueva instancia de CatController.
	catController := NewCatController()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer catController.DbConnection.Close()

	// Obtiene la colección de gatos.
	collection := catController.DbConnection.GetCollection(CollectionNameCat)
	// Variable que contiene a todos los gatos
	var cats []models.Cat
	// Trae a todos los gatos desde la base de datos
	results, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Println("No fue posible traer a todos los gatos")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for results.Next(ctx) {
		var singleCat models.Cat
		if err = results.Decode(&singleCat); err != nil {
			log.Println("Usuario no se pudo añadir")
		}

		cats = append(cats, singleCat)
	}
	ctx.JSON(http.StatusOK, cats)
}
