package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func DoctorRoutes(app *fiber.App) {
	app.Get("/doctors", controllers.GetAllDoctors)
	app.Get("/doctors/available", controllers.GetAvailableDoctors)
	// app.Get("/doctors/:name", controllers.GetDoctorByName)
	app.Get("/doctors/:id", controllers.GetDoctorById)
	app.Get("/doctors/:doctorId/slots", controllers.GetDoctorScheduleById)
}
