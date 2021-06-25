package services

import (
	dataaccess "tinc1/DataAccess"
	models "tinc1/Models"
)

type FilesService interface {
	GetInboundFiles(id string) []models.Inboundfile
}

type filesService struct {
	dao dataaccess.Dao
}

func DBFilesService(dao dataaccess.Dao) FilesService {
	return &filesService{
		dao: dao,
	}
}

func (service *filesService) GetInboundFiles(id string) []models.Inboundfile {
	data := service.dao.GetInboundFiles(id)
	return data
}
