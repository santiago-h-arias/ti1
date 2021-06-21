package dto

type LoginCredentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Files struct {
	Id string `form: "id", binding:"required"`
}
