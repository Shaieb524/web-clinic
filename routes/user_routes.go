package routes

import (
	"github.com/Shaieb524/web-clinic.git/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/user/:userId", controllers.GetAUser)
	router.GET("/users", controllers.GetAllUsers)
}
