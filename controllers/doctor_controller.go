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

func GetAllDoctors(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	query := bson.M{"role": "doctor"}
	opts := options.Find().SetProjection(bson.D{{"schedule", 0}})

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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"doctors": users}},
	)
}

func GetDoctorByName(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nameParam := c.Params("name")
	query := bson.D{{Key: "name", Value: nameParam}}

	var doctorDoc bson.M
	var err error
	if err = userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println(err)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"doctor": doctorDoc}},
	)
}

func GetDoctorById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	idParam := c.Params("id")
	doctorId, err := primitive.ObjectIDFromHex(idParam)
	query := bson.D{{Key: "_id", Value: doctorId}}

	var doctorDoc bson.M
	if err = userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println(err)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"doctor": doctorDoc}},
	)
}

func GetDoctorScheduleById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	idParam := c.Params("doctorId")
	doctorId, err := primitive.ObjectIDFromHex(idParam)

	query := bson.D{{Key: "_id", Value: doctorId}}

	var doctorDoc bson.M
	if err = userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println(err)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"schedule": doctorDoc["schedule"]}},
	)
}

func GetAvailableDoctors(c *fiber.Ctx) error {
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
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"available_doctors": users}},
	)
}
