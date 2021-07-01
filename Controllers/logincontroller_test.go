package controllers

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"regexp"
	"testing"
	dto "tinc1/Dto"
	services "tinc1/Services"
	testutils "tinc1/TestUtils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupLoginController() (mock sqlmock.Sqlmock, fitted_loginController LoginController, err error) {
	//MAKE JWTSERVICE
	mocked_jwtService := testutils.NewMock_JWTService()

	//MAKE DAO
	mock_db, mock, err := sqlmock.New() //Mock Instance of sql.DB
	if err != nil {
		return
	}

	//sqlX.DB instance with core as a mocked sql.DB
	mock_xdb := sqlx.NewDb(mock_db, "sqlserver")

	//Our dataAccess operator, uses sqlX.DB
	mock_dao := testutils.NewMock_Dao(*mock_xdb)

	fitted_loginService := services.DBLoginService(mock_dao)
	fitted_loginController = NewLoginController(fitted_loginService, mocked_jwtService)

	return
}

func Test_Login(t *testing.T) {
	mock, fitted_loginController, err := setupLoginController()
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)

	t.Run("validCredentials", func(t *testing.T) {
		creds := dto.LoginCredentials{
			Email:    "ajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		//Preparing a Request Body
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

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

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
		token, _ := fitted_loginController.Login(mock_ctx)

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		//Output should be untampered
		if token != "1234_valid_token_4321" {
			t.Fatalf("\nCouldn't authenticate\nExpected : \"%s\"\tGot : \"%s\"", "1234_valid_token_4321", token)
		}
	})

	t.Run("invalidCredentials", func(t *testing.T) {
		creds := dto.LoginCredentials{
			Email:    "notajith@thinkbridge.in",
			Password: "Ajith12#",
		}

		//Preparing a Request Body
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

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"NaesbUserKey",
				"Name",
				"Email",
			}) //Empty (Now Rows)
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		//Execute funcion to be tested
		token, _ := fitted_loginController.Login(mock_ctx)

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		//Output should be untampered
		if token != "" {
			t.Fatalf("\nShouldn't have authenticated!\nExpected : \"%s\"\tGot : \"%s\"", "", token)
		}
	})

	t.Run(("noCredentials"), func(t *testing.T) {
		creds := dto.LoginCredentials{
			Email:    "",
			Password: "",
		}

		//Preparing a Request Body
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

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"NaesbUserKey",
				"Name",
				"Email",
			})
		mock.ExpectQuery(regexp.QuoteMeta("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2")).WillReturnRows(rows)

		//Execute funcion to be tested
		token, _ := fitted_loginController.Login(mock_ctx)

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		//Output should be untampered
		if token != "" {
			t.Fatalf("\nShouldn't have authenticated!\nExpected : \"%s\"\tGot : \"%s\"", "", token)
		}
	})

	t.Run(("emptyRequest"), func(t *testing.T) {
		//Preparing a Context containing (Empty) Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/login", nil)

		//Execute funcion to be tested
		token, _ := fitted_loginController.Login(mock_ctx)

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}

		//Output should be untampered
		if token != "" {
			t.Fatalf("\nShouldn't have authenticated!\nExpected : \"%s\"\tGot : \"%s\"", "", token)
		}
	})
}
