package models

import "database/sql"

type Outboundfile struct {
	OutboundFileKey           string         `db:"OutboundFileKey" json:"outbound_file_key,omitempty"`
	NaesbUserKey              string         `db:"NaesbUserKey" json:"naesb_user_key,omitempty"`
	UsKey                     string         `db:"UsKey" json:"us_key,omitempty"`
	UsCommonCode              string         `db:"UsCommonCode" json:"us_common_code,omitempty"`
	ThemKey                   string         `db:"ThemKey" json:"them_key,omitempty"`
	ThemCommonCode            string         `db:"ThemCommonCode" json:"them_common_code,omitempty"`
	Filename                  string         `db:"Filename" json:"file_name,omitempty"`
	Plaintext                 string         `db:"Plaintext" json:"plain_text,omitempty"`
	Ciphertext                sql.NullString `db:"Ciphertext" json:"-"`
	Attempt1At                sql.NullString `db:"Attempt1At" json:"attempt_1_at,omitempty"`
	Attempt2At                sql.NullString `db:"Attempt2At" json:"attempt_2_at,omitempty"`
	Attempt3At                sql.NullString `db:"Attempt3At" json:"attempt_3_at,omitempty"`
	Receipt                   sql.NullString `db:"Receipt" json:"receipt,omitempty"`
	Result                    sql.NullString `db:"Result" json:"result,omitempty"`
	Escalated                 string         `db:"Escalated" json:"escalated,omitempty"`
	EscalatedAt               sql.NullString `db:"EscalatedAt" json:"escalated_at,omitempty"`
	Debug                     string         `db:"Debug" json:"debug,omitempty"`
	CreatedAt                 string         `db:"CreatedAt" json:"created_at,omitempty"`
	EmpowerOutgoingEdiFileKey sql.NullString `db:"EmpowerOutgoingEdiFileKey" json:"empower_outgoing_edi_file_key,omitempty"`
	DoNotSend                 string         `db:"DoNotSend" json:"do_not_send,omitempty"`
	LastLocation              sql.NullString `db:"LastLocation" json:"last_location,omitempty"`
	Ciphered                  string         `db:"Ciphered" json:"-"`
	Posted                    string         `db:"Posted" json:"posted,omitempty"`
	OutboundFileId            string         `db:"OutboundFileId" json:"outbound_file_id,omitempty"`
	Inactive                  string         `db:"Inactive" json:"inactive,omitempty"`
}
