package services

import (
	dataaccess "tinc1/DataAccess"
	models "tinc1/Models"
)

type LoginService interface {
	Login(username string, password string) (bool, models.NaesbUser)
}

type loginService struct {
	dao dataaccess.Dao
}

func DBLoginService(dao dataaccess.Dao) LoginService {
	return &loginService{
		dao: dao,
	}
}

func (service *loginService) Login(email string, password string) (bool, models.NaesbUser) {
	isAuthenticated, user := service.dao.CheckUser(email, password)
	return isAuthenticated, user
}
