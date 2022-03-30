package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/models"
	"github.com/Shaieb524/web-clinic.git/responses"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func GetAUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		fmt.Println("error getting a user : ", userId, user)
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
}

func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		users = append(users, singleUser)
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}},
	)
}
