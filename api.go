package main

// TODO: Refactoring. Move general purpose functions into a new utilities package
// TODO: Testing.

import (
	"net/http"
	"os"

	controllers "tinc1/Controllers"
	dataaccess "tinc1/DataAccess"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

var (
	// TODO: implement files service
	// TODO: implement files dao
	// filesService  services.FileService = service.NewFileService()
	// filesDao        dataaccess.Dao      = dataaccess.NewDao()
	dao        dataaccess.Dao      = dataaccess.NewDao()
	jwtService services.JWTService = services.NewJWTService()

	// TODO: implement files controller
	// filesController controllers.filesController = controllers.NewFilesController()
	loginService    services.LoginService       = services.DBLoginService(dao)
	loginController controllers.LoginController = controllers.NewLoginController(loginService, jwtService)
)

func main() {

	api := gin.New()

	api.Use(gin.Recovery(), gin.Logger())

	// Login Endpoint: Authentication + Token creation + User response
	api.POST("/login", func(ctx *gin.Context) {
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	api.Run(":" + port)
}
