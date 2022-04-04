package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Shaieb524/web-clinic.git/models"
	"github.com/Shaieb524/web-clinic.git/responses"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo"
)

func GetAllPatients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	query := bson.M{"role": "patient"}
	opts := options.Find().SetProjection(bson.D{{"_id", 0}})

	results, err := userCollection.Find(ctx, query, opts)

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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"patients": users}},
	)
}

func GetPatientById(c *gin.Context) {
	idParam := c.Param("id")
	patientDoc := getDPatientProfileByStringId(idParam)

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"patient": patientDoc}},
	)
}

func getDPatientProfileByStringId(strPatientId string) bson.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	patientId, err := primitive.ObjectIDFromHex(strPatientId)
	if err != nil {
		fmt.Println("error converting id from hex")
	}

	query := bson.D{{Key: "_id", Value: patientId}}

	var patientDoc bson.M
	if err := userCollection.FindOne(ctx, query).Decode(&patientDoc); err != nil {
		fmt.Println(err)
	}

	return patientDoc
}
