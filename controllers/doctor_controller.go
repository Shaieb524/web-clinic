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

func GetAllDoctors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	query := bson.M{"role": "doctor"}
	opts := options.Find().SetProjection(bson.D{{"schedule", 0}})

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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"doctors": users}},
	)
}

func GetDoctorByName(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nameParam := c.Param("name")
	query := bson.D{{Key: "name", Value: nameParam}}

	var doctorDoc bson.M
	var err error
	if err = userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"doctor": doctorDoc}},
	)
}

func GetDoctorById(c *gin.Context) {
	idParam := c.Param("id")
	doctorDoc := getDoctorProfileByStringId(idParam)

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"doctor": doctorDoc}},
	)
}

func GetDoctorScheduleById(c *gin.Context) {
	idParam := c.Param("id")
	doctorDoc := getDoctorProfileByStringId(idParam)

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"schedule": doctorDoc["schedule"]}},
	)
}

func getDoctorProfileByStringId(strDoctorId string) bson.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doctorId, err := primitive.ObjectIDFromHex(strDoctorId)
	if err != nil {
		fmt.Println("error converting id from hex")
	}

	query := bson.D{{Key: "_id", Value: doctorId}}

	var doctorDoc bson.M
	if err := userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println(err)
	}

	return doctorDoc
}

func GetAvailableDoctors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	query := bson.M{
		"$and": []bson.M{
			{"role": "doctor"},
			{"available": true},
		},
	}
	opts := options.Find().SetProjection(bson.D{{"schedule", 0}})

	results, err := userCollection.Find(ctx, query, opts)
	fmt.Println("available results : ", results)
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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"available_doctors": users}},
	)
}
