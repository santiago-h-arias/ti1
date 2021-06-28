package models

type Outboundfile struct {
	OutboundFileKey           string `db:"OutboundFileKey" json:"File_key"`
	UsKey                     string `db:"UsKey" json:"user_key"`
	UsCommonCode              string `db:"UsCommonCode" json:"UsCommonCode"`
	ThemKey                   string `db:"ThemKey" json:"them_key"`
	ThemCommonCode            string `db:"ThemCommonCode" json:"ThemCommonCode"`
	Filename                  string `db:"Filename" json:"Filename"`
	Plaintext                 string `db:"Plaintext" json:"Plaintext"`
	Ciphertext                string `db:"Ciphertext" json:"Ciphertext"`
	Attempt1At                string `db:"Attempt1At" json:"file_attempt1at"`
	Attempt2At                string `db:"Attempt2At" json:"file_attempt2at"`
	Attempt3At                string `db:"Attempt3At" json:"file_attempt3at"`
	Receipt                   string `db:"Receipt" json:"file_receipt"`
	Result                    string `db:"Result" json:"file_result"`
	Escalated                 string `db:"Escalated" json:"file_escalated"`
	EscalatedAt               string `db:"EscalatedAt" json:"file_escalated_at"`
	Debug                     string `db:"Debug" json:"file_debug"`
	CreatedAt                 string `db:"CreatedAt" json:"file_created_at"`
	EmpowerOutgoingEdiFileKey string `db:"EmpowerOutgoingEdiFileKey" json:"file_empower_key"`
	DoNotSend                 string `db:"DoNotSend" json:"file_do_not_send"`
	LastLocation              string `db:"LastLocation" json:"file_last_location"`
	Ciphered                  string `db:"Ciphered" json:"file_ciphered"`
	Posted                    string `db:"Posted" json:"file_posted"`
	OutboundFileId            string `db:"OutboundFileId" json:"outbound_id"`
	NaesbUserKey              string `db:"NaesbUserKey" json:"NaesbUserKey"`
	Inactive                  string `db:"Inactive" json:"Inactive"`
}
