package controller

/*
import (
	"context"
	"net/http"
	"rest-template/config"
	"rest-template/models"
	"rest-template/responses"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionNameCat = "Cat"
)

var catCollection *mongo.Collection = config.GetCollection(config.DB, CollectionNameCat)

func CreateCat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var cat models.Cat
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&cat); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		/*
			//use the validator library to validate required fields
			if validationErr := validate.Struct(&user); validationErr != nil {
				c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
				return
			}
*/
/*
		newUser := models.Cat{
			ID:        primitive.NewObjectID(),
			Name:      cat.Name,
			Breed:     cat.Breed,
			Age:       cat.Age,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(), // Revisar la coma final

		}

		result, err := catCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
*/
