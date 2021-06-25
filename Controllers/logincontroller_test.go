package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"regexp"
	"testing"
	dataaccess "tinc1/DataAccess"
	dto "tinc1/Dto"
	models "tinc1/Models"
	services "tinc1/Services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

//For Mocking JWTService Interface
type mock_jwtService struct {
	secretkey string
	issuer    string
}

func (jwtService *mock_jwtService) GenerateToken(email string) string {
	if email == "ajith@thinkbridge.in" {
		return "1234_valid_token_4321"
	}
	fmt.Printf("\n\nExpected:%s\nGot:%s\n\n\n", "ajith@thinkbridge.in", email)

	return ""
}

func (jwtService *mock_jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return &jwt.Token{}, nil
}

func NewMock_JWTService() services.JWTService {
	return &mock_jwtService{
		secretkey: "UjgFm344XW",
		issuer:    "thinkbridgeIdProvider",
	}
}

//For Mocking Dao Interface (for LoginService)
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

//For Mocking io.Reader Interface (for gin.Context)
type Reader struct {
	data      []byte
	readIndex int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(len(r.data)) {
		err = io.EOF
		return
	}

	n = copy(p, r.data[r.readIndex:])
	r.readIndex += int64(n)
	return
}

func Test_Login(t *testing.T) {
	//MAKE JWTSERVICE
	mocked_jwtService := NewMock_JWTService()

	//MAKE DAO
	mock_db, mock, err := sqlmock.New() //Mock Instance of sql.DB
	if err != nil {
		fmt.Println("expected no error, but got:", err)
		return
	}
	defer mock_db.Close()

	mock_xdb := sqlx.NewDb(mock_db, "sqlserver") //sqlX.DB instance with core as a mocked sql.DB
	mock_dao := NewMock_Dao(*mock_xdb)           //Our dataAccess operator, uses sqlX.DB

	//MAKE LOGIN SERVICE USING DAO
	mocked_loginService := services.DBLoginService(mock_dao)

	//MAKE LOGINCONTROLLER using LOGINSERVICE and JWTSERVICE
	fitted_loginController := NewLoginController(mocked_loginService, mocked_jwtService)

	gin.SetMode(gin.TestMode)

	t.Run("validCredentials", func(t *testing.T) {
		creds := dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("email", creds.Email))
		assert.NoError(t, mw.WriteField("password", creds.Password))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Emulated output
		rows := sqlmock.NewRows([]string{"NaesbUserKey", "Name", "Email"}).
			AddRow("8B0528AB-6E22-40E2-9B60-A4A6C584E6E3", "Ajith", "ajith@thinkbridge.in")
		//Query to expect and then emulate output
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		token, _ := fitted_loginController.Login(mock_ctx)

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		if token != "1234_valid_token_4321" {
			t.Fatalf("\nCouldn't authenticate\nExpected : \"%s\"\tGot : \"%s\"", "1234_valid_token_4321", token)
		}
	})

	t.Run("invalidCredentials", func(t *testing.T) {
		creds := dto.LoginCredentials{
			Email:    "notajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("email", creds.Email))
		assert.NoError(t, mw.WriteField("password", creds.Password))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Emulated output
		rows := sqlmock.NewRows([]string{"NaesbUserKey", "Name", "Email"})
		//Query to expect and then emulate output
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		token, _ := fitted_loginController.Login(mock_ctx)

		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		if token != "" {
			t.Fatalf("\nShouldn't have authenticated!\nExpected : \"%s\"\tGot : \"%s\"", "", token)
		}
	})

}
