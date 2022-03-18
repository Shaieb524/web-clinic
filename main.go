package main

import (
	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()

	app.Listen(":6000")
}
