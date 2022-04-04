package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"
	"github.com/gin-gonic/gin"
)

func PatientRoutes(router *gin.Engine) {
	router.GET("/patients", controllers.GetAllPatients)
	router.GET("/patients/:id", controllers.GetDoctorById)
	router.GET("/patients/:id/appointments-history", controllers.ViewPatientAppointmentsHistory)
}
