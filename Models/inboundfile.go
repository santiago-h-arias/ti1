package models

type Inboundfile struct {
	InboundFileKey string `db:"InboundFileKey" json:"file_key"`
	UsKey          string `db:"UsKey" json:"user_key"`
	UsCommonCode   string `db:"UsCommonCode" json:"us_common_code"`
	ThemKey        string `db:"ThemKey" json:"them_key"`
	ThemCommonCode string `db:"ThemCommonCode" json:"them_common_code"`
	Filename       string `db:"Filename" json:"file_name"`
	Plaintext      string `db:"Plaintext" json:"file_plain_text"`
	Ciphertext     string `db:"Ciphertext" json:"-"`
	ReceivedAt     string `db:"ReceivedAt" json:"file_received_at"`
	TransactionId  string `db:"TransactionId" json:"transaction_id"`
	Processed      string `db:"Processed" json:"processed"`
	InboundFileId  string `db:"InboundFileId" json:"inbound_file_id"`
	NaesbUserKey   string `db:"NaesbUserKey" json:"naesb_user_key"`
	Inactive       string `db:"Inactive" json:"inactive"`
}
