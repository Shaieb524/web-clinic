package main

import (
	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()

	// unauthorized routes
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "pong"})
	})
	routes.UnauthRoute(app)

	// restricted routes
	routes.UserRoute(app)

	app.Listen(":" + configs.EnvPort())
}
