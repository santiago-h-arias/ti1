package main

// TODO: Refactoring. Move general purpose functions into a new utilities package
// TODO: Testing.

import (
	"net/http"
	"os"

	controllers "tinc1/Controllers"
	dataaccess "tinc1/DataAccess"
	middlewares "tinc1/Middlewares"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

var (
	dao        dataaccess.Dao      = dataaccess.NewDao()
	jwtService services.JWTService = services.NewJWTService()

	filesService    services.FilesService       = services.DBFilesService(dao)
	filesController controllers.FilesController = controllers.NewFilesController(filesService)
	loginService    services.LoginService       = services.DBLoginService(dao)
	loginController controllers.LoginController = controllers.NewLoginController(loginService, jwtService)
)

func main() {

	api := gin.New()

	api.Use(gin.Recovery(), gin.Logger(), CORS)

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

	// Protectected routes are grouped here
	apiRoutes := api.Group("/api", middlewares.AuthorizeJWT())
	{
		apiRoutes.POST("/files", func(c *gin.Context) {
			data := filesController.GetInboundFiles(c)
			c.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	api.Run(":" + port)
}

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
