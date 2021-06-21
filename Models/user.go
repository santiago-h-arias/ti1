package models

import "github.com/google/uuid"

type NaesbUser struct {
	Naesbuserkey string
	Name         string
	Email        string
	Isadmin      bool
}
type Inboundfile struct {
	InboundFileKey string
	UsKey          uuid.UUID `gorm:"type:uuid;" json:"user_key"`
	UsCommonCode   string
	ThemKey        uuid.UUID `gorm:"type:uuid;" json:"theme_key"`
	ThemCommonCode string
	Filename       string
	Plaintext      string
	CipherText     string
	ReceivedAt     string
	TransactionId  string
	Processed      string
	InboundFileId  string
}
