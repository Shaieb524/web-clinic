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
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

func GetAllDoctors(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{"role": "doctor"})

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
	fmt.Println("gegt doc by name")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// doctorName := c.Params("doctorName")

	// result, err := userCollection.FindOne(ctx, bson.M{"name": "sawsan"})

	// var user models.User
	// res := userCollection.FindOne(ctx, bson.M{"name": "sawsan"}).Decode(&user)?

	filter := bson.M{"name": "sawsan"}
	singleResult := userCollection.FindOne(ctx, filter)
	fmt.Println("sing res : ", singleResult)

	// if err != nil {
	// 	fmt.Println("err 1 : ", err)
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	// }
	// fmt.Println("err 2 : ", err)

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"doctor": singleResult}},
	)
}

func GetDoctorById(c *fiber.Ctx) int {
	fmt.Println("gegt doc by id")
	return 1
	// get id by params
	// params := c.Params("id")

	// _id, err := primitive.ObjectIDFromHex(params)
	// if err != nil {
	// 	return c.Status(500).SendString(err.Error())
	// }

	// filter := bson.D{{"_id", _id}}

	// var result models.User

	// if err != userCollection.FindOne(c.Context(), filter).Decode(&result); err != nil {
	// 	return c.Status(500).SendString("Something went wrong.")
	// }

	// return c.Status(fiber.StatusOK).JSON(result)
}
