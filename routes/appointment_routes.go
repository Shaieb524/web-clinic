package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func AppointmentRoutes(app *fiber.App) {
	app.Post("/appointments/book-appointment", controllers.BookAppointmentSlot)
}
