package routes

import (
	"jwtapi/controllers"

	"github.com/gin-gonic/gin"
)

func Paths(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/profile", controllers.Profile)
	r.POST("/logout", controllers.Logout)
}
