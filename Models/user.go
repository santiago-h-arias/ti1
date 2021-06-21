package models

import "github.com/google/uuid"

type NaesbUser struct {
	NaesbUserKey uuid.UUID `gorm:"type:string;primary_key" json:"user_key"`
	Name         string    `gorm:"type:string" json:"user_name"`
	Email        string    `gorm:"type:string" json:"user_email"`
	IsAdmin      bool      `gorm:"type:bool" json:"user_admin"`
}
