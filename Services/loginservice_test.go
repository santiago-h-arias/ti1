package services

import (
	"fmt"
	"regexp"
	"testing"
	dto "tinc1/Dto"
	testutils "tinc1/TestUtils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func setupLoginService() (mock sqlmock.Sqlmock, fitted_LoginService LoginService, err error) {
	//Mock Instance of sql.DB
	mock_db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)
		return
	}

	//sqlX.DB instance with core as a mocked sql.DB
	mock_xdb := sqlx.NewDb(mock_db, "sqlserver")

	//Our dataAccess operator, uses sqlX.DB
	mock_dao := testutils.NewMock_Dao(*mock_xdb)
	fitted_LoginService = DBLoginService(mock_dao)

	return
}

func Test_Login(t *testing.T) {
	mock, fitted_LoginService, err := setupLoginService()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("validCredentials", func(t *testing.T) {
		//Sample Valid Creds, Doesn't really matter because Query output is hardcoded
		creds := &dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"NaesbUserKey",
				"Name",
				"Email",
			}).
			AddRow("8B0528AB-6E22-40E2-9B60-A4A6C584E6E3", "Ajith", "ajith@thinkbridge.in")
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		//Execute funcion to be tested
		auth, _ := fitted_LoginService.Login(creds.Email, creds.Password)

		//Should Authenticate
		if !auth {
			t.Fatalf(`Failed to Authenticate`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})

	t.Run("invalidCredentials", func(t *testing.T) {
		//Sample Invalid Creds, Doesn't really matter because Query output is hardcoded
		creds := &dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "notAjith12#",
		}

		//Execute funcion to be tested
		auth, _ := fitted_LoginService.Login(creds.Email, creds.Password)

		//Should Authenticate
		if auth {
			t.Fatalf(`Shouldn't have authenticated!`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})
}
