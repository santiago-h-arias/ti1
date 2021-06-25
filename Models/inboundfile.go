package models

type Inboundfile struct {
	InboundFileKey string `db:"InboundFileKey" json:"File_key"`
	UsKey          string `db:"UsKey" json:"user_key"`
	UsCommonCode   string `db:"UsCommonCode" json:"UsCommonCode"`
	ThemKey        string `db:"ThemKey" json:"them_key"`
	ThemCommonCode string `db:"ThemCommonCode" json:"ThemCommonCode"`
	Filename       string `db:"Filename" json:"Filename"`
	Plaintext      string `db:"Plaintext" json:"Plaintext"`
	Ciphertext     string `db:"Ciphertext" json:"Ciphertext"`
	ReceivedAt     string `db:"ReceivedAt" json:"ReceivedAt"`
	TransactionId  string `db:"TransactionId" json:"TransactionId"`
	Processed      string `db:"Processed" json:"Processed"`
	InboundFileId  string `db:"InboundFileId" json:"InboundFileId"`
	NaesbUserKey   string `db:"NaesbUserKey" json:"NaesbUserKey"`
	Inactive       string `db:"Inactive" json:"Inactive"`
}
