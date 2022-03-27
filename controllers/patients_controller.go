package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Shaieb524/web-clinic.git/models"
	"github.com/Shaieb524/web-clinic.git/responses"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo"
)

func GetAllPatients(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	query := bson.M{"role": "patient"}
	opts := options.Find().SetProjection(bson.D{{"_id", 0}})

	results, err := userCollection.Find(ctx, query, opts)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"patients": users}},
	)
}

func GetPatientById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	patientDoc := getDPatientProfileByStringId(idParam)

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"patient": patientDoc}},
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
