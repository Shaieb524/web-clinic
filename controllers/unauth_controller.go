package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/helpers"
	"github.com/Shaieb524/web-clinic.git/models"
	"github.com/Shaieb524/web-clinic.git/responses"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Ping(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{"message": "pong"})
}

func RegisterUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//encypt user's password before saving to db
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	newUser := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		Role:     data["role"],
	}

	// if user is a doctor then initialize a schedule
	if data["role"] == "doctor" {
		ds := models.DoctorSchedule(helpers.GenerateWeekDoctorSchedule("insetredStrId"))
		newUser.Schedule = ds
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"user": result}})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": data["email"]}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"message": "Invalid Credentials"}})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"message": "Invalid Credentials"}})
	}

	pair, err := generateTokenPair(user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"message": "Login failed!"}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"pair ": pair}})
}

type tokenPair struct {
	Token        string
	RefreshToken string
}

func generateTokenPair(userEmail string) (tokenPair, error) {

	var tPair tokenPair

	// create token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userEmail,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})

	token, err := claims.SignedString([]byte(configs.EnvTokenSecretKey()))
	if err != nil {
		return tPair, err
	}

	// create refresh token
	rtclaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userEmail,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	rtoken, err := rtclaims.SignedString([]byte(configs.EnvRefreshTokenSecretKey()))
	if err != nil {
		return tPair, err
	}

	tPair = tokenPair{Token: token, RefreshToken: rtoken}

	return tPair, nil
}

// type handler struct{}

// This is the api to refresh tokens
// Most of the code is taken from the jwt-go package's sample codes
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
// func RefreshToken(c echo.Context) error {
// 	type tokenReqBody struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}
// 	tokenReq := tokenReqBody{}
// 	c.Bind(&tokenReq)

// 	// Parse takes the token string and a function for looking up the key.
// 	// The latter is especially useful if you use multiple keys for your application.
// 	// The standard is to use 'kid' in the head of the token to identify
// 	// which key to use, but the parsed token (head and claims) is provided
// 	// to the callback, providing flexibility.
// 	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 		return []byte("secret"), nil
// 	})

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// Get the user record from database or
// 		// run through your business logic to verify if the user can log in
// 		if int(claims["sub"].(float64)) == 1 {

// 			newTokenPair, err := generateTokenPair("joud@gmail.com")
// 			if err != nil {
// 				return err
// 			}

// 			return c.JSON(http.StatusOK, newTokenPair)
// 		}

// 		return echo.ErrUnauthorized
// 	}

// 	return err
// }
