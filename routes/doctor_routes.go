package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func DoctorRoutes(app *fiber.App) {
	app.Get("/doctors", controllers.GetAllDoctors)
	// app.Get("/doctors/:doctorName", controllers.GetDoctorByName)
	// app.Get("/doctors/:doctorId", controllers.GetDoctorById)
}
