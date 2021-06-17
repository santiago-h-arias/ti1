package main

import (
	"net/http"
	"os"

	controllers "tinc1/Controllers"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

var (
	// filesService  services.FileService = service.NewFileService()
	loginService services.LoginService = services.DBLoginService()
	jwtService   services.JWTService   = services.NewJWTService()

	// filesController controllers.filesController = controllers.NewFilesController()
	loginController controllers.LoginController = controllers.NewLoginController(loginService, jwtService)
)

func main() {

	api := gin.New()

	api.Use(gin.Recovery(), gin.Logger())

	// Login Endpoint: Authentication + Token creation
	api.POST("/login", func(ctx *gin.Context) {
		// fmt.Println(dao.CheckUser("ajith@thinkbridge.in", "Ajith12#"))
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	// To set the port as an env variable from the eb console
	port := os.Getenv("PORT")
	// Elastic Beanstalk forwards requests to port 8000
	if port == "" {
		port = "8000"
	}
	api.Run(":" + port)
}
