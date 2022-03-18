package main

import (
	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "pong"})
	})

	configs.ConnectDB()
	routes.UserRoute(app)

	app.Listen(":6000")
}
