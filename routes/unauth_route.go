package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func UnauthRoute(app *fiber.App) {
	app.Post("/user/register", controllers.RegisterUser)
	app.Post("/user/login", controllers.Login)
}
