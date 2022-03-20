package main

import (
	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()

	// logger middleware
	app.Use(logger.New())

	// unauthorized routes
	routes.UnauthRoutes(app)

	// JWT validation
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(configs.EnvSecretKey()),
	}))

	// restricted routes
	routes.UserRoutes(app)

	app.Listen(":" + configs.EnvPort())
}
