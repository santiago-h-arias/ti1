package controllers

import (
	dto "tinc1/Dto"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	loginService services.LoginService
	jWtService   services.JWTService
}

func NewLoginController(loginService services.LoginService,
	jWtService services.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	var credentials dto.LoginCredentials
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		return ""
	}
	isAuthenticated := controller.loginService.Login(credentials.Email, credentials.Password)
	if isAuthenticated {
		return controller.jWtService.GenerateToken(credentials.Email)
	}
	return ""
}
