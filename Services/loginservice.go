package services

// dao          dataaccess.Dao        = dataaccess.NewDao()
type LoginService interface {
	Login(username string, password string) bool
}

type loginService struct {
	email    string
	password string
}

func DBLoginService() LoginService {
	return &loginService{
		email:    "test@thinkbridge.com",
		password: "test",
	}
}

func (service *loginService) Login(email string, password string) bool {
	return service.email == email && service.password == password
}
