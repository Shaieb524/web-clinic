package main

import (
	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/routes"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()

	// unauthorized routes
	routes.UnauthRoutes(app)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(configs.EnvSecretKey()),
	}))

	// restricted routes
	routes.UserRoutes(app)

	app.Listen(":" + configs.EnvPort())
}
