package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Colecciones en la base de datos
type DbConnection struct {
	Client  *mongo.Client
	Context context.Context
}

// Función que estable conexión a la base de datos
func NewDbConnection() *DbConnection {
	//Url para la conexión a mongodb
	uri := "mongodb://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_URL") + "/" + os.Getenv("DB_DB")
	log.Println(uri)
	//Se establece la conexión con la base de datos
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbNames)

	return &DbConnection{
		Client:  client,
		Context: ctx,
	}
}

// Función que cierra la conexión a la base de datos
func (c *DbConnection) Close() {
	c.Client.Disconnect(c.Context)
}

// Función que retorna una colección
func (c *DbConnection) GetCollection(collection string) *mongo.Collection {
	return c.Client.Database(os.Getenv("DB_DB")).Collection(collection)
}
