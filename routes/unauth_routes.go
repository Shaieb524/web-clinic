package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func UnauthRoutes(app *fiber.App) {
	app.Get("/ping", controllers.Ping)
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.Login)
	// app.Post("/refresh-token", controllers.)
}
