package main

import (
	// "fmt"
	// "time"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/routes"

	"github.com/gin-gonic/gin"
	// 	jwt "github.com/appleboy/gin-jwt/v2"
)

func main() {
	app := gin.Default()
	configs.ConnectDB()

	// unauthorized routes
	routes.UnauthRoutes(app)
	// restricted routes
	routes.UserRoutes(app)
	routes.DoctorRoutes(app)
	routes.PatientRoutes(app)
	routes.AppointmentRoutes(app)

	app.Run(":" + configs.EnvPort())
}
