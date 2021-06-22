package main

import (
	"fmt"
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

	filesdao dataaccess.Files = dataaccess.NewFiles()

	fileservice     services.FileService              = services.DBFilesService(filesdao)
	filesController controllers.FileServiceController = controllers.NewFilesController(fileservice)
	loginService    services.LoginService             = services.DBLoginService(dao)
	loginController controllers.LoginController       = controllers.NewLoginController(loginService, jwtService)
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

	api.POST("/files", func(c *gin.Context) {
		data := filesController.Files(c)
		fmt.Printf("data: %s", data)
		// err := c.ShouldBindJSON(&data)
		// if err != nil {
		// 	panic(err)
		// }
		// if data !=  {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"data": data,
		// 	})
		// }
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})
	// To set the port as an env variable from the eb console
	port := os.Getenv("PORT")
	// Elastic Beanstalk forwards requests to port 8000
	if port == "" {
		port = "8000"
	}
	api.Run(":" + port)
}
