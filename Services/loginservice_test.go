package services

import (
	"fmt"
	"regexp"
	"testing"
	dataaccess "tinc1/DataAccess"
	dto "tinc1/Dto"
	models "tinc1/Models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type mock_dao struct {
	db sqlx.DB
}

func (dao *mock_dao) CheckUser(email string, password string) (bool, models.NaesbUser) {
	var user models.NaesbUser

	err := dao.db.QueryRowx("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2", email, password).StructScan(&user)
	if err == nil {
		return true, user
	}
	return false, user

}

func NewMock_Dao(db sqlx.DB) dataaccess.Dao {
	return &mock_dao{
		db: db,
	}
}

func (dao *mock_dao) GetInboundFiles(id string) []models.Inboundfile {
	return []models.Inboundfile{}
}

func (dao *mock_dao) GetOutboundFiles(id string) []models.Outboundfile {
	return []models.Outboundfile{}
}

func Test_Login(t *testing.T) {
	//Mock Instance of sql.DB
	mock_db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)
		return
	}
	defer mock_db.Close()

	//sqlX.DB instance with core as a mocked sql.DB
	mock_xdb := sqlx.NewDb(mock_db, "sqlserver")

	//Emulated output
	rows := sqlmock.NewRows([]string{"NaesbUserKey", "Name", "Email"}).
		AddRow("8B0528AB-6E22-40E2-9B60-A4A6C584E6E3", "Ajith", "ajith@thinkbridge.in")
	//Query to expect and then emulate output
	mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

	//Our dataAccess operator, uses sqlX.DB
	mock_dao := NewMock_Dao(*mock_xdb)
	mock_LoginService := DBLoginService(mock_dao)

	t.Run("validCredentials", func(t *testing.T) {
		//Sample Valid Creds
		creds := &dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "Ajith12#",
		}
		auth, _ := mock_LoginService.Login(creds.Email, creds.Password)

		if !auth {
			t.Fatalf(`Failed to Authenticate`)
		}

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})

	t.Run("invalidCredentials", func(t *testing.T) {
		//Sample Invalid Creds
		creds := &dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "notAjith12#",
		}
		auth, _ := mock_LoginService.Login(creds.Email, creds.Password)

		if auth {
			t.Fatalf(`Shouldn't have authenticated!`)
		}

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})
}
