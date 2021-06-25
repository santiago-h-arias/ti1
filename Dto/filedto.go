package dto

type FileRequestDto struct {
	Id string `json:"id" binding:"required"`
}
