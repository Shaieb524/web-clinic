package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"
	"github.com/gin-gonic/gin"
)

func AppointmentRoutes(router *gin.Engine) {
	router.POST("/appointments/book-appointment", controllers.BookAppointmentSlot)
	router.POST("/appointments/cancel-appointment", controllers.CancelAppointmentSlot)
	router.POST("/appointments/view-appointment", controllers.ViewAppointmentDetails)
}
