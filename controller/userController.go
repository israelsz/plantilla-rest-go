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
	"golang.org/x/crypto/bcrypt"
)

/*
Se establecen los nombres de la colección que se traeran desde la base de datos
*/
const (
	CollectionNameUser = "User"
)

// Función para crear un usuario e insertarlo a la base de datos de mongodb
func CreateUser(ctx *gin.Context) {
	log.Println("Service: CreateUser")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()

	// Obtiene la colección de usuarios.
	collection := dbConnection.GetCollection(CollectionNameUser)

	// Obtiene los datos del usuario a partir del cuerpo de la solicitud HTTP.
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Println("La estructura ingresada no es válida")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Genera un nuevo ID único para el usuario.
	user.ID = primitive.NewObjectID()
	// Establece la fecha de creación y actualización del usuario.
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// Encriptación de contraseña
	userPass := ctx.Param("password")
	user.Hash = GeneratePassword(userPass)

	// Inserta el usuario en la colección.
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("No fue posible crear un Usuario")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Devuelve el usuario creado.
	ctx.JSON(http.StatusCreated, user)
	log.Println("Usuario creado")
}

// Función para obtener un gato por id
func GetUserByID(ctx *gin.Context) {
	log.Println("Service: GetUserByID")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	userID := ctx.Param("id")

	// Crea un objeto ID de MongoDB a partir del ID del gato.
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("No fue posible convertir el ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"_id": oid}

	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameCat)

	// Ejecuta el método FindOne() con el filtro y el contexto.
	var result models.User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
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
	log.Println("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	ctx.JSON(http.StatusOK, result)
}

// Función para obtener un gato por id
func GetUserByEmail(ctx *gin.Context) {
	log.Println("Service: GetUserByID")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	userEmail := ctx.Param("email")
	log.Println("User email:", userEmail)
	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"email": userEmail}

	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameUser)

	// Ejecuta el método FindOne() con el filtro y el contexto.
	var result models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			log.Println("Usuario no encontrado")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		// Ocurrió un error durante la búsqueda.
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	ctx.JSON(http.StatusOK, result)
}

// Funciones de autenticación

func ComparePasswords(storedHash string, loginPass string) error {
	byteHash := []byte(storedHash)
	loginHash := []byte(loginPass)
	err := bcrypt.CompareHashAndPassword(byteHash, loginHash)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GeneratePassword(password string) string {
	binpwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(binpwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
