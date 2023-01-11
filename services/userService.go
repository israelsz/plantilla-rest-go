package services

import (
	"log"
	"rest-template/config"
	"rest-template/models"
	"rest-template/utils"
	"time"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Se establecen los nombres de la colección que se traeran desde la base de datos
*/
const (
	CollectionNameUser = "User"
)

// Función para crear un usuario e insertarlo a la base de datos de mongodb
func CreateUserService(newUser models.User) (models.User, error) {
	log.Println("Service: CreateUser")
	//Se valida el usuario antes de ingresar a la base de datos
	ok, err := govalidator.ValidateStruct(newUser)
	log.Println("Valor de ok:", ok)
	//Si el usuario no tiene una estructura valida
	if !ok {
		log.Println("Validation error: ", err)
		return newUser, err
	}
	//Si el usuario es valido
	//Se establece conexión con base de datos mongo
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene la colección de usuarios.
	collection := dbConnection.GetCollection(CollectionNameUser)

	// Se revisa si el usuario se encuentra en la base de datos

	// Buscar si el email existe
	var result models.User
	// Creando un filtro de busqueda
	filter := bson.M{"email": newUser.Email}
	err = collection.FindOne(dbConnection.Context, filter).Decode(&result)
	//Si no fue encontrar el email
	if err != nil {
		//Si el email no se encuentra en la base de datos
		if err == mongo.ErrNoDocuments {
			newUser.ID = primitive.NewObjectID()
			// Establece la fecha de creación y actualización del gato.
			newUser.CreatedAt = time.Now()
			newUser.UpdatedAt = time.Now()
			// Se encripta la contraseña
			newUser.Hash = utils.GeneratePassword(newUser.Password)
			// Se vacia el campo password
			newUser.Password = ""
			// No se encontró ningún documento con el email especificado, entonces se inserta el nuevo usuario
			_, err = collection.InsertOne(dbConnection.Context, newUser)
			if err != nil {
				log.Println("Error al insertar nuevo usuario: ", err)
				return newUser, err
			}
			log.Println("Nuevo usuario creado con exito")
			return newUser, nil
		}
		// Ocurrió un error durante la búsqueda.
		log.Println("Error al buscar el email en la base de datos: ", err)
		return newUser, err
	}
	log.Println("El email ya existe en la base de datos:", err)
	return newUser, err
}

// Función para obtener un gato por id
func GetUserByIDService(userID string) (models.User, error) {
	log.Println("Service: GetUserByID")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	// Crea un objeto ID de MongoDB a partir del ID del gato.
	var result models.User
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("No fue posible convertir el ID")
		return result, err
	}
	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"_id": oid}

	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameUser)
	err = collection.FindOne(dbConnection.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			log.Println("Usuario no encontrado")
			return result, err
		}
		// Ocurrió un error durante la búsqueda.
		return result, err
	}
	log.Println("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	return result, nil
}

// Función para obtener un gato por id
func GetUserByEmailService(userEmail string) (models.User, error) {
	log.Println("Service: GetUserByEmail")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	var result models.User
	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"email": userEmail}

	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameUser)
	err := collection.FindOne(dbConnection.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el email especificado.
			log.Println("Usuario no encontrado, err")
			return result, err
		}
		// Ocurrió un error durante la búsqueda.
		return result, err
	}
	log.Println("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	return result, nil
}

func UpdateUserService(updatedUser models.User, userID string) (models.User, error) {
	log.Println("Service: UpdateUser")
	//Se valida el usuario antes de ingresar a la base de datos
	ok, err := govalidator.ValidateStruct(updatedUser)
	//Si el usuario no tiene una estructura valida
	if !ok {
		log.Println("Validation error: ", err)
		return updatedUser, err
	}
	var resultUser models.User
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("No fue posible convertir el ID")
		return resultUser, err
	}
	// Se actualiza la fecha de actualización
	resultUser.UpdatedAt = time.Now()
	update := bson.M{"$set": updatedUser}
	filter := bson.M{"_id": oid}
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameUser)
	_, err = collection.UpdateOne(dbConnection.Context, filter, update)
	if err != nil {
		return resultUser, err
	}
	log.Println("Usuario actualizado")
	return resultUser, nil
}
