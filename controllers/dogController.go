package controllers

import (
	"context"
	"net/http"
	"time"
)

func createDog(c *gin.Context) {
	var dog Dog
	c.BindJSON(&dog)
	dog.ID = primitive.NewObjectID()
	dog.CreatedAt = time.Now()
	dog.UpdatedAt = time.Now()
	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	dogsCollection := db.(*mongo.Client).Database("dogs").Collection("dogs")
	_, err := dogsCollection.InsertOne(context.TODO(), dog)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       dog.ID,
		"message":  "Dog created successfully!",
		"resource": dog,
	})
}
