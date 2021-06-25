package controllers

import (
	dto "tinc1/Dto"
	models "tinc1/Models"
	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

type FilesController interface {
	GetInboundFiles(ctx *gin.Context) []models.Inboundfile
}

type filesController struct {
	fileService services.FilesService
}

func NewFilesController(fileService services.FilesService) FilesController {
	return &filesController{
		fileService: fileService,
	}
}

func (controller *filesController) GetInboundFiles(ctx *gin.Context) []models.Inboundfile {
	var fileId dto.FileRequestDto
	err := ctx.ShouldBind(&fileId)
	if err != nil {
		return []models.Inboundfile{}
	}
	data := controller.fileService.GetInboundFiles(fileId.Id)

	// return []models.Inboundfile{}
	return data

}
