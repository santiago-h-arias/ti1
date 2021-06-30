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

func NewMock_Dao(db sqlx.DB) dataaccess.Dao {
	return &mock_dao{
		db: db,
	}
}

func (dao *mock_dao) CheckUser(email string, password string) (bool, models.NaesbUser) {
	var user models.NaesbUser

	err := dao.db.QueryRowx("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2", email, password).StructScan(&user)
	if err == nil {
		return true, user
	}
	return false, user

}

func (dao *mock_dao) GetInboundFiles(id string) []models.Inboundfile {
	var findfiles []models.Inboundfile
	err := dao.db.Select(&findfiles, "select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(InboundFileKey as char(36)) as InboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from InboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey  where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1", id)
	if err != nil {
		fmt.Println(err)
	}
	return findfiles
}

func (dao *mock_dao) GetOutboundFiles(id string) []models.Outboundfile {
	var findfiles []models.Outboundfile
	err := dao.db.Select(&findfiles, "select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(OutboundFileKey as char(36)) as OutboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from OutboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1", id)
	if err != nil {
		fmt.Println(err)
	}
	return findfiles
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
