package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"

	"github.com/gofiber/fiber/v2"
)

func PatientRoutes(app *fiber.App) {
	app.Get("/patients", controllers.GetAllPatients)
	app.Get("/patients/:id", controllers.GetDoctorById)
	app.Get("/patients/:id/appointments-history", controllers.ViewPatientAppointmentsHistory)
}
