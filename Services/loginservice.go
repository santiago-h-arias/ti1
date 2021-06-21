package services

import (
	"fmt"
	dataaccess "tinc1/DataAccess"
	models "tinc1/Models"
)

type LoginService interface {
	Login(username string, password string) (bool, models.NaesbUser)
}

type FileService interface {
	Files(id string) []models.Inboundfile
}

type loginService struct {
	dao dataaccess.Dao
}

type filesService struct {
	dao dataaccess.Files
}

func DBLoginService(dao dataaccess.Dao) LoginService {
	return &loginService{
		dao: dao,
	}
}

func DBFilesService(dao dataaccess.Files) FileService {
	return &filesService{
		dao: dao,
	}
}

func (service *loginService) Login(email string, password string) (bool, models.NaesbUser) {
	isAuthenticated, user := service.dao.CheckUser(email, password)
	return isAuthenticated, user
}

func (service *filesService) Files(id string) []models.Inboundfile {
	data := service.dao.Get_files(id)
	fmt.Printf("service:%s", data)
	return data
}
