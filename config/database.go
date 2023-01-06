package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadDatabase() {
	//Url para la conexión a mongodb
	uri := "mongodb://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_URL") + "/" + os.Getenv("DB_DB")
	fmt.Println(uri)
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
	defer client.Disconnect(ctx)

	// Sacar esto despues de ver si funciono la conexion
	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//databases, err := client.ListDatabaseNames(ctx, bson.M{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(databases)

}

//os.Getenv("DB_URL")

//mongodb://templategoRESTUser:Sf17a033vcF!@localhost:27017/?authMechanism=DEFAULT
//mongodb://os.Getenv("DB_USER"):os.Getenv("DB_PASS")@os.Getenv("DB_URL")/?authMechanism=DEFAULT

//clientOptions := options.Client().ApplyURI(middlewares.DotEnvVariable("MONGO_URL"))
//mongodb://mongo:27100/db

//https://www.mongodb.com/blog/post/quick-start-golang-mongodb-starting-and-setup

// mongodb://templategoRESTUser:Sf17a033vcF!@localhost:27017/templategoREST
