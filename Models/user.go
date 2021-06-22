package models

type NaesbUser struct {
	User_key   string `db:"NaesbUserKey" json:"user_key"`
	User_name  string `db:"Name" json:"user_name"`
	User_email string `db:"Email" json:"user_email"`
	User_admin bool   `db:"IsAdmin" json:"user_admin"`
}
