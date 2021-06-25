package dataaccess

import (
	"fmt"
	"regexp"
	"testing"
	"tinc1/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

//Test authentication with valid Credentials
func Test_CheckUser(t *testing.T) {
	//Mock Instance of sql.DB
	mock_db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)
		return
	}
	defer mock_db.Close()

	//sqlX.DB instance with core as a mocked sql.DB
	mock_xdb := sqlx.NewDb(mock_db, "sqlserver")

	//Our dataAccess operator, uses sqlX.DB
	mock_dao := &dao{db: mock_xdb}
	//Call for credential match into database

	t.Run("validCredentials", func(t *testing.T) {
		//Sample Valid Creds
		creds := &dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		//Emulated output
		rows := sqlmock.NewRows([]string{"NaesbUserKey", "Name", "Email"}).
			AddRow("8B0528AB-6E22-40E2-9B60-A4A6C584E6E3", "Ajith", "ajith@thinkbridge.in")
		//Query to expect and then emulate output
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		isAuthenticated, _ := mock_dao.CheckUser(creds.Email, creds.Password)

		//Should Authenticate
		if !isAuthenticated {
			t.Fatalf(`Failed to Authenticate`)
		}

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})

	t.Run("invalidCredentials", func(t *testing.T) {
		//Sample Inalid Creds
		creds := &dto.LoginCredentials{
			Email:    "notajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		//Emulated output
		rows := sqlmock.NewRows([]string{"NaesbUserKey", "Name", "Email"})
		//Query to expect and then emulate output
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		isAuthenticated, _ := mock_dao.CheckUser(creds.Email, creds.Password)

		//Should not authenticate!
		if isAuthenticated {
			t.Fatalf(`Shouldn't have authenticated`)
		}

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})
}
