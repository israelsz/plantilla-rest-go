package services

import (
	"errors"
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

// Función que valida un usuario
func validarUsuarioCreate(user models.User) (bool, error) {
	err := errors.New("Estructura invalida")
	// Se verifica que el email ingresado sea valido
	if !govalidator.IsEmail(user.Email) {

		return false, err
	}
	// Se verifica que haya un nombre ingresado
	if !govalidator.ByteLength(user.Name, "2", "20") || !govalidator.IsAlphanumeric(user.Name) {
		return false, err
	}
	// Se verifica que haya una contraseña ingresada
	if !govalidator.ByteLength(user.Password, "2", "20") || !govalidator.IsAlphanumeric(user.Password) {
		return false, err
	}
	// Se verifica que el rol sea ingresado
	if !govalidator.ByteLength(user.Rol, "1", "20") {
		return false, err
	}
	return true, nil
}

// Función que valida un usuario
func validarUsuarioUpdate(user models.User) (bool, error) {
	err := errors.New("Estructura invalida")
	// Se verifica que el email ingresado sea valido
	if govalidator.IsEmail(user.Email) && !(user.Email == "") {
		err = correoOcupado(user.Email)
		if err != nil {
			return false, err
		}

	}
	// Se verifica que haya un nombre ingresado
	if (!govalidator.ByteLength(user.Name, "2", "20") || govalidator.IsAlphanumeric(user.Name)) && user.Name != "" {
		return false, err
	}
	// Se verifica que haya una contraseña ingresada
	if (!govalidator.ByteLength(user.Password, "2", "20") || !govalidator.IsAlphanumeric(user.Password)) && user.Password != "" {
		return false, err
	}
	// Se verifica que el rol sea ingresado
	if !govalidator.ByteLength(user.Rol, "1", "20") && user.Rol != "" {
		return false, err
	}
	return true, nil
}

func correoOcupado(email string) error {
	dbConnection := config.NewDbConnection()
	defer dbConnection.Close()
	var result models.User
	filter := bson.M{"email": email}
	collection := dbConnection.GetCollection(CollectionNameUser)
	err := collection.FindOne(dbConnection.Context, filter).Decode(&result)
	return err
}

// Función para crear un usuario e insertarlo a la base de datos de mongodb
func CreateUserService(newUser models.User) (models.User, error) {
	utils.Debug("Service: CreateUser")
	//Se valida el usuario antes de ingresar a la base de datos
	ok, err := validarUsuarioCreate(newUser)
	utils.Debug("Valor de ok:", ok)
	//Si el usuario no tiene una estructura valida
	if !ok {
		utils.Debug("Estructura invalida")
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
	err = correoOcupado(newUser.Email)
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
				utils.Debug("Error al insertar nuevo usuario: ", err)
				return newUser, err
			}
			utils.Debug("Nuevo usuario creado con exito")
			return newUser, nil
		}
		// Ocurrió un error durante la búsqueda.
		return newUser, err
	}
	return newUser, errors.New("usuario se encuentra en la base de datos")
}

// Función para obtener a un usuario por su id
func GetUserByIDService(userID string) (models.User, error) {
	utils.Debug("Service: GetUserByID")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Crea un objeto ID de MongoDB a partir del ID del usuario
	var result models.User
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.Debug("No fue posible convertir el ID")
		return result, err
	}
	// Crea un filtro para buscar al usuario por su ID.
	filter := bson.M{"_id": oid}
	// Obtiene la colección de usuarios.
	collection := dbConnection.GetCollection(CollectionNameUser)
	err = collection.FindOne(dbConnection.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el ID especificado.
			utils.Debug("Usuario no encontrado")
			return result, err
		}
		// Ocurrió un error durante la búsqueda.
		return result, err
	}
	utils.Debug("Se encontró el usuario")
	// Devuelve al usuario encontrado.
	return result, nil
}

// Función para obtener a un usuario por id
func GetUserByEmailService(userEmail string) (models.User, error) {
	utils.Debug("Service: GetUserByEmail")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	var result models.User
	// Crea un filtro para buscar al usuario por su ID.
	filter := bson.M{"email": userEmail}

	// Obtiene la colección de usuarios.
	collection := dbConnection.GetCollection(CollectionNameUser)
	err := collection.FindOne(dbConnection.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No se encontró ningún documento con el email especificado.
			utils.Debug("Usuario no encontrado, err")
			return result, err
		}
		// Ocurrió un error durante la búsqueda.
		return result, err
	}
	utils.Debug("Se encontró el usuario")
	// Devuelve el usuario encontrado.
	return result, nil
}

func GetAllUserService() ([]models.User, error) {
	utils.Debug("Service: GetAllUserService")
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	collection := dbConnection.GetCollection(CollectionNameUser)
	// Variable que contiene a todos los usuarios en un arreglo
	var users []models.User
	// Trae a todos los usuarios desde la base de datos
	results, err := collection.Find(dbConnection.Context, bson.M{})
	if err != nil {
		return users, errors.New("no fue posible traer a todos los usuarios")
	}
	for results.Next(dbConnection.Context) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			utils.Debug("Usuario no se pudo añadir")
		}

		users = append(users, singleUser)
	}
	return users, nil
}

func UpdateUserService(updatedUser models.User, userID string) (models.User, error) {
	utils.Debug("Service: UpdateUser")
	ok, err := validarUsuarioUpdate(updatedUser)
	//Se valida el usuario antes de ingresar a la base de datos
	//Si el usuario no tiene una estructura valida
	if !ok {
		utils.Debug("Validation error: ", err)
		return updatedUser, err
	}
	var resultUser models.User
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.Debug("No fue posible convertir el ID")
		return resultUser, err
	}
	if updatedUser.Password != "" {
		updatedUser.Hash = utils.GeneratePassword(updatedUser.Password)
		updatedUser.Password = ""
	}
	// Se actualiza la fecha de actualización
	resultUser.UpdatedAt = time.Now()
	update := bson.M{"$set": updatedUser}
	filter := bson.M{"_id": oid}
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Obtiene la colección de usuarios.
	collection := dbConnection.GetCollection(CollectionNameUser)
	_, err = collection.UpdateOne(dbConnection.Context, filter, update)
	if err != nil {
		return resultUser, err
	}
	utils.Debug("Usuario actualizado")
	return resultUser, nil
}

func DeleteUserService(userID string) error {
	// Crea un objeto ID de MongoDB a partir del ID del usuario.
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.Debug("No fue posible convertir el ID")
		return errors.New("id invalido")
	}
	// Crea una nueva instancia a la conexión de base de datos
	dbConnection := config.NewDbConnection()
	// Define un defer para cerrar la conexión a la base de datos al finalizar la función.
	defer dbConnection.Close()
	// Filtro para la query
	filter := bson.M{"_id": oid}
	collection := dbConnection.GetCollection(CollectionNameUser)
	// Elimina al usuario de la colección.
	result, _ := collection.DeleteOne(dbConnection.Context, filter)
	utils.Debug(result)
	// Si no hay error
	if result.DeletedCount == 1 {
		// Se pudo eliminar el usuario
		return nil
	}
	// No se pudo eliminar el usuario
	return errors.New("usuario no pudo ser eliminado")
}
