package services

import (
	"errors"
	"log"
	"rest-template/config"
	"rest-template/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Se establecen los nombres de la colección que se traeran desde la base de datos
*/
const (
	CollectionNameCat = "Cat"
)

// Función para crear un gato e insertarlo a la base de datos de mongodb
func CreateCatService(newCat models.Cat) (models.Cat, error) {
	//Se establece conexión con la base de datos mongo
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameCat)

	// Genera un nuevo ID único para el gato.
	newCat.ID = primitive.NewObjectID()
	// Establece la fecha de creación y actualización del gato.
	newCat.CreatedAt = time.Now()
	newCat.UpdatedAt = time.Now()

	// Inserta el gato en la colección.
	_, err := collection.InsertOne(dbConnection.Context, newCat)
	//Si hubo un error
	if err != nil {
		return newCat, err
	}

	return newCat, err
}

// Función para obtener un gato por id
func GetCatByIDService(catID string) (models.Cat, error) {
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Crea un objeto ID de MongoDB a partir del ID del gato.
	var gato models.Cat
	oid, err := primitive.ObjectIDFromHex(catID)
	if err != nil {
		log.Println("No fue posible convertir el ID")
		return gato, err
	}
	// Crea un filtro para buscar el gato por su ID.
	filter := bson.M{"_id": oid}

	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameCat)
	err = collection.FindOne(dbConnection.Context, filter).Decode(&gato)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			log.Println("Gato no encontrado")
			return gato, err
		}
		// Ocurrió un error durante la búsqueda.
		return gato, err
	}
	log.Println("Se encontró el gato")
	// Devuelve el gato encontrado.
	return gato, nil

}

func DeleteCatService(catID string) error {
	// Obtiene el ID del gato a partir del parámetro de la ruta.
	// Crea un objeto ID de MongoDB a partir del ID del gato.
	oid, err := primitive.ObjectIDFromHex(catID)
	if err != nil {
		log.Println("No fue posible convertir el ID")
		return errors.New("id invalido")
	}
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Se elimina el gato
	filter := bson.M{"_id": oid}
	collection := dbConnection.GetCollection(CollectionNameCat)
	// Elimina el gato de la colección.
	result, _ := collection.DeleteOne(dbConnection.Context, filter)
	// Si no hay error
	if result.DeletedCount == 1 {
		// Se pudo eliminar al gato
		return nil
	}
	// No se pudo eliminar al gato
	return errors.New("el gato no pudo ser eliminado")
}

func UpdateCatService(updatedCat models.Cat, userID string) (models.Cat, error) {
	// Se crea el modelo del gato que se actualizara
	var resultCat models.Cat
	//Se crea el objectId de mongo
	oid, err := primitive.ObjectIDFromHex(userID)
	//Si ocurrio un error
	if err != nil {
		log.Println("No fue posible convertir el ID")
		return resultCat, err
	}
	// Se actualiza la fecha de actualización
	resultCat.UpdatedAt = time.Now()
	update := bson.M{"$set": updatedCat}
	filter := bson.M{"_id": oid}
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene la colección de gatos.
	collection := dbConnection.GetCollection(CollectionNameCat)
	_, err = collection.UpdateOne(dbConnection.Context, filter, update)
	if err != nil {
		return resultCat, err
	}
	log.Println("Gato actualizado")
	return resultCat, nil
}

func GetAllCatsService() ([]models.Cat, error) {
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	collection := dbConnection.GetCollection(CollectionNameCat)
	// Variable que contiene a todos los gatos
	var cats []models.Cat
	// Trae a todos los gatos desde la base de datos
	results, err := collection.Find(dbConnection.Context, bson.M{})
	if err != nil {
		return cats, errors.New("no fue posible traer a todos los gatos")
	}
	for results.Next(dbConnection.Context) {
		var singleCat models.Cat
		if err = results.Decode(&singleCat); err != nil {
			log.Println("Gato no se pudo añadir")
		}

		cats = append(cats, singleCat)
	}
	return cats, nil
}
