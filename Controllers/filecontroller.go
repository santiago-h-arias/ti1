package controllers

import (
	"fmt"
	dto "tinc1/Dto"
	models "tinc1/Models"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

type FileServiceController interface {
	Files(ctx *gin.Context) []models.Inboundfile
}

type filesController struct {
	fileService services.FileService
}

func NewFilesController(fileService services.FileService) FileServiceController {
	return &filesController{
		fileService: fileService,
	}
}

func (controller *filesController) Files(ctx *gin.Context) []models.Inboundfile {
	var files dto.Files
	err := ctx.ShouldBind(&files)
	if err != nil {
		return []models.Inboundfile{}
	}
	data := controller.fileService.Files(files.Id)
	fmt.Printf("hey:%s", data)

	// return []models.Inboundfile{}
	return data

}
