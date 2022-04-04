package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"
	"github.com/gin-gonic/gin"
)

func DoctorRoutes(router *gin.Engine) {
	router.GET("/doctors", controllers.GetAllDoctors)
	router.GET("/doctors/available", controllers.GetAvailableDoctors)
	// router.GET("/doctors/:name", controllers.GetDoctorByName)
	router.GET("/doctors/:id", controllers.GetDoctorById)
	router.GET("/doctors/:id/slots", controllers.GetDoctorScheduleById)
}
