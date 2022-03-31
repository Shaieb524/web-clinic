package main

import (
	// "fmt"
	// "time"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/middlewares"
	"github.com/Shaieb524/web-clinic.git/routes"

	"github.com/gin-gonic/gin"
	// 	jwt "github.com/appleboy/gin-jwt/v2"
)

func main() {
	router := gin.Default()
	configs.ConnectDB()

	// unauthorized routes
	routes.UnauthRoutes(router)

	// jwt middleware auth
	router.Use(middlewares.JWTauthentication())

	// restricted routes
	routes.UserRoutes(router)
	routes.DoctorRoutes(router)
	routes.PatientRoutes(router)
	routes.AppointmentRoutes(router)

	router.Run(":" + configs.EnvPort())
}
