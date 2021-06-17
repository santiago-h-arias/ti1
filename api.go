package main

import (
	"net/http"
	"os"

	controllers "tinc1/Controllers"
	dataaccess "tinc1/DataAccess"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

var (
	// filesService  services.FileService = service.NewFileService()
	dao        dataaccess.Dao      = dataaccess.NewDao()
	jwtService services.JWTService = services.NewJWTService()

	// filesController controllers.filesController = controllers.NewFilesController()
	loginService    services.LoginService       = services.DBLoginService(dao)
	loginController controllers.LoginController = controllers.NewLoginController(loginService, jwtService)
)

func main() {

	api := gin.New()

	api.Use(gin.Recovery(), gin.Logger())

	// Login Endpoint: Authentication + Token creation
	api.POST("/login", func(ctx *gin.Context) {
		// fmt.Println(dao.CheckUser("ajith@thinkbridge.in", "Ajith123#"))
		token, user := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
				"user":  user,
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
