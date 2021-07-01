package controllers

import (
	"bytes"
	"database/sql"
	"mime/multipart"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	models "tinc1/Models"
	services "tinc1/Services"
	testutils "tinc1/TestUtils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupfilesController() (mock sqlmock.Sqlmock, fitted_filesController FilesController, err error) {
	//Mock Instance of sql.DB
	mock_db, mock, err := sqlmock.New()
	if err != nil {
		return
	}

	//sqlX.DB instance with core as a mocked sql.DB
	mock_xdb := sqlx.NewDb(mock_db, "sqlserver")

	//Our dataAccess operator, uses sqlX.DB
	mock_dao := testutils.NewMock_Dao(*mock_xdb)

	fitted_filesService := services.DBFilesService(mock_dao)
	fitted_filesController = NewFilesController(fitted_filesService)

	return
}

func TestGetInboundFiles(t *testing.T) {
	mock, fitted_filesController, err := setupfilesController()
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)

	t.Run("validFileId", func(t *testing.T) {
		sampleId := "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3" //Doesn't really matter because Query output is hardcoded

		//Preparing a Request Body
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("Id", sampleId))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/api/inboundfiles", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Sample DB output, can be anything
		sample_row1 := models.Inboundfile{
			InboundFileKey: "F3F53EB3-91DD-484E-B38F-0063D03343A8",
			UsKey:          "D6DFAD93-8B80-4119-8B94-6814F1ED75BD",
			UsCommonCode:   "829244552110T",
			ThemKey:        "00999A26-0D64-414D-949D-78ADA131AC46",
			ThemCommonCode: "183529049",
			Filename:       "867FWD.446407.1.edi",
			Plaintext:      "ISA*00*          *00*          *01*183529049      *01*829244552110T  *161028*0732*U*00401*000000015*0*T*>~GS*PT*183529049*8292445521100*20161028*073228*1*X*004010~ST*867*1001~BPT*SU*CPLA737883B119523136*20161028~REF*Q5**10032789419418520~REF*TN*EPSD814T16T1610271103C4~N1*8S*AEP TEXAS CENTRAL*1*007924772~N1*AY*ERCOT*1*183529049**41~N1*SJ*ALLTEX ALLTEX POWER AND LIGHT*9*8292445521100**40~PTD*BJ***MG*134404448~DTM*140*20161027~QTY*QD***NV~MEA****KH**0*51~SE*12*1001~GE*1*1~IEA*1*000000015~",
			Ciphertext:     "--=--\r\nContent-Type: application/pgp-encrypted\r\n\r\nVersion: 1\r\n\r\n--=--\r\nContent-Type: application/octet-stream\r\n\r\n-----BEGIN PGP MESSAGE-----\r\nVersion: PGP SDK 3.0.3\r\n\r\nqANQR1DBwU4D/X9rcfytvJgQB/wOO0eNEt91BE0wYhhBWmT27GRe2A15vzgFADrC\r\nIcysQmbmq20qiBVZ53T88acO2iE0pBgAG8S9A3PYhoPsH1absWo/Hp03JVRUOLXK\r\nANKVsb7zl3Fy+GcDyVPH4rDG1Gf5G0ngWc4zKj7bfGyZXo5yp8AfQpb7FQcHLpMi\r\nwU7bu7h9ypalwejOh/IFLPWih2HlOTun9EyStSfXWkoQcSw1xEgwNo1Yawx9EezM\r\nSoppCXcNejsUs/u769bKcoS7bx/EqiYXZPE7824Y141IYjL9zz4/ztYX5AjHFhrT\r\nIXvSNAOMA9EkX8z0uC3KrJpBhC+iVBJSHv0ptfnj8JsC4QBxB/4rpKJlxxf++z9/\r\nryRql41iU0aQ33eQiGAtXlcKVZC39psHmBYxZ7ckL2r139QDocCvrjMBwTU+6DLe\r\nqr5R208VHdViYAj5Uei8UcIsq2DZWjwLEAQYO4yMhzTwcm8BoAUlCfkAa8/tWakX\r\n8Jb26BrzgkXDZBwo1VHxtC2RpN8S30GLSr6lVECbxnArvdq/Jm8zBZ0Cwwk6RZGb\r\nGDZz9u/+k477WFAEbybuo+0QaBsYFNu4uzIr2JMbzng6MqUyV0kzCsP+Hji74G+W\r\nbm8Y7PwcYkzKYToZBL7BnZi/ckNehUCorks8tU/ev7SExvZep6aXxP7YFsr4ABUl\r\n3r4xg0gX0sGIAfjjhNWV7yywYY1v6cT8q2Y+Zz51P0GjMV517kG7uQ8b0YbpiB9v\r\nAYPD/uHq5zmunVMrNjyTVwQbMlkT87cF+VCvCfTpauKCvZgwQHEzBjfytD3kL+0B\r\nfNjs03flmY0TTF8J84ZSYLo9M0Xgla+AuIT4NuqwsFD+QdKUIQNEAHs2kt0aQ7pQ\r\nkJ0U+MkcdXX/Nd21rgqAjROilo/Tr8VKjRP1r5GcBTHaFzvNX+06N3SNOxDWMNTa\r\nGwleYpaW7lVq5ESejxBOdDLLtPNbi0DWvrQXxHmdZmePiGWuLVfvshS558M0wY5h\r\noTtiVo/tLINnlGzB2f8HQOtVuMJDMdCLrMVC0GsEMMpxOsOD+DxtL9x9LSODIjtk\r\nnowIRjuEK1o89YTWwTpgreou7vz+eUM4FFH218uMLquizT5Lu3KggQvKILhbxhPs\r\nDEQLnxpN2FRjqGtmSQQLPmJ2gR4QmFDSPMwNbh6jzBqZsUnT9+gJsC/M9F/oIQ1h\r\nlmeXRzTUzSXZfchOebmrTNvg2ii/o3omUBPGesBqtPdn4iHY+7uOimy22uS7QDcv\r\nRlhHAaD7c1z1ABSDo/9/UofoK4iN7gA0xCGrC2Gf+yIQ7jpVfyZKJw3SKdBBNlwU\r\nRuLyLQXLNO+j0pYzk54qwLRHAxKuCjB4yaC/QXMe0nMnRzu/mcuiVyKD9IApoFJi\r\nwF6aqt0eikXR8CCmqyPpNQHFF6K4xXICpwO/ihZ1MtPalP0tqgwLesDb03NlufGL\r\nfqi5IIT/Db4tmAYLIN5kI9c=\r\n=kfFj\r\n-----END PGP MESSAGE-----\r\n\r\n--=----\r\n",
			ReceivedAt:     "2016-10-28T07:34:03.343Z",
			TransactionId:  "161028073403081",
			Processed:      "true",
			InboundFileId:  "1",
			NaesbUserKey:   "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3",
			Inactive:       "false",
		}
		sample_row2 := models.Inboundfile{
			InboundFileKey: "F6837013-CAE7-44EE-8D57-01201A46AF55",
			UsKey:          "D6DFAD93-8B80-4119-8B94-6814F1ED75BD",
			UsCommonCode:   "829244552110T",
			ThemKey:        "B6B5856B-D996-45AF-B63A-E868B1BF5497",
			ThemCommonCode: "1052623364500",
			Filename:       "D:\\CORE\\Temp\\SHXTX82924455A_ILLUMINARV2_CONNECT_81",
			Plaintext:      "ISA*00*          *00*          *01*829244552110T  *01*1052623364500T *160916*0937*U*00401*000066779*0*T*>~GS*GE*8292445521100*1052623364500*20160916*0937175*66877*X*004010~ST*814*170850~BGN*13*2016091609334602*20160916~N1*8S*SHARYLAND*1*1052623364500~N1*SJ*Illuminar Energy*1*8292445521100~N1*8R*Lord Vader~LIN*2*SH*EL*SH*CE~ASI*7*021~REF*11*0977777777~REF*12*0977777777~REF*BLT*DUAL~REF*PC*DUAL~REF*9V*N~SE*13*170850~GE*1*66877~IEA*1*000066779~",
			Ciphertext:     "------------ECP153039\r\nContent-Type: application/pgp-encrypted\r\n\r\nVersion: 1\r\n\r\n------------ECP153039\r\nContent-Type: application/octet-stream\r\n\r\n-----BEGIN PGP MESSAGE-----\r\nVersion: PGP Command Line v10.4.0 (Build 94) (Win32)\r\n\r\nqANQR1DBwU4D/X9rcfytvJgQB/0Q0hwGaQFF+J5ShTAfYXS+qIjiXg6quPxZvpqB\r\n5NQvMAtJkv7Qhuf+0sqPxHGxmzCruX1iMJBeePwxXtx+0n4Qh3c2OZEClNChkMdu\r\nlzrI0cMPIefI721RobImYAH74ky88sw5GYftKKXyeq7vJV58cFExo91Cf1QDq2vz\r\nzroLHiWb4UfYeKPKfH+33LT8XFxwFcsyuJBjMh2HSaOKCyhiWJR2QwGE+gJRq7qx\r\n0fnrAwX9d266z483wKg2Oy8aWLhYwfuqUxzlE4JgICBVEqqBhJj2Ev5n5nMdeFMP\r\ndaSY/VRHbavsc5573S4Z7Nfkp3faZn5JksSYaum/YQ/sw0UtB/46MNGwU6b4DUm/\r\nHEMHwBFaeIapRNZZ0Cl2o3xRNrAKagmEnOgg/oF+rfu+Gdg6cQicp+v8soPIZ+Xb\r\nzciyWL/tfE1KgnKD/D2yf0prWO/vJrbi+gMfyRdSu5M9x9skgYpN+hYBpHW7KwdL\r\nFWFa7kWBYcHg8BvTzfl1hOpgXjFd3N4XgNtfvzW2ecEHkwUgz+d0CM6l77UAChrf\r\nw32HgaH6OPDdlhi1dAu1ceQpHTiTKB92Jn3321Z5k7CqsIB89KfD/tCoqJZ4cCfo\r\nkTymRMXsJC1Xzue6Ebe76BDD3++HNJ2VnBsixfQRBEFJzYsQzzl4zIIGlzrpiLY6\r\nNyNPvSYV0sEBAYcIC+Q/74mDFSvURGoeKsEIapUHSbkJ5PXW6PZ0HKPHk/nTVUu7\r\nrerhARvdpaZpXXmZcli041xbasLcdUSZLKGPRaGm0QDiKKzrYrFhy5+jo2yxO+QQ\r\nJ7JrwRw08QYzIKu/mU7eKvh9csYULAM0HVPVpucvharoa2aoXnLpdm2Cxj5PCae4\r\nJmH4mqEl5ZokVg/QzVuiaaCo/vmHim4JRW7bIvXfH4H0bnUwpS5gwuL0Tknp3254\r\n/jtdtRxHyk8NfNDuCtE7IQmLRXGxw9EtaKTQKjIovRI5NztAfs1MRX1HW3Y81PXa\r\n/aptwokfzVt3Smx3B08VXD4D9MUTYbcmNDjDoTIskjo7plSSWbW3fMXEtHw7YlGW\r\nEoju/ehRa4w+28o5RDMgqYBrch+me4xLN4mwN6PBqacZS/Rml4Tyqd7Sy39Kagxr\r\nKRfQ6grpxOG6+288ea0fn0QYZBMAv7mi9Rlx07KO4ApVuvF3n9iICEpYUQ6frDWt\r\nCU3/ano7XYHBJWcvZ2dFil6Pe6eVBApDaxCK49wNIA1tA0VPOKrNFxS1fe0JuAAA\r\nBEShlJjzRenDD69GT3pA8iqFotCOykJb/Kg=\r\n=UXgC\r\n-----END PGP MESSAGE-----\r\n\r\n------------ECP153039--",
			ReceivedAt:     "2016-10-24T15:31:12.19Z",
			TransactionId:  "161024153112055",
			Processed:      "true",
			InboundFileId:  "4",
			NaesbUserKey:   "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3",
			Inactive:       "false",
		}

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"InboundFileKey",
				"UsKey",
				"UsCommonCode",
				"ThemKey",
				"ThemCommonCode",
				"Filename",
				"Plaintext",
				"Ciphertext",
				"ReceivedAt",
				"TransactionId",
				"Processed",
				"InboundFileId",
				"NaesbUserKey",
				"Inactive",
			}).
			AddRow(testutils.RowFromStruct(sample_row1)...).
			AddRow(testutils.RowFromStruct(sample_row2)...)
		mock.ExpectQuery(regexp.QuoteMeta("select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(InboundFileKey as char(36)) as InboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from InboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey  where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1")).WillReturnRows(rows)

		//Execute funcion to be tested
		output := fitted_filesController.GetInboundFiles(mock_ctx)

		//Output should be untampered
		if !reflect.DeepEqual(output, []models.Inboundfile{sample_row1, sample_row2}) {
			t.Fatalf(`Output doesn't match`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})

	t.Run("invalidFileId", func(t *testing.T) {
		sampleId := "69420" //Doesn't really matter because Query output is hardcoded

		//Preparing a Request Body
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("Id", sampleId))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/api/inboundfiles", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"InboundFileKey",
				"UsKey",
				"UsCommonCode",
				"ThemKey",
				"ThemCommonCode",
				"Filename",
				"Plaintext",
				"Ciphertext",
				"ReceivedAt",
				"TransactionId",
				"Processed",
				"InboundFileId",
				"NaesbUserKey",
				"Inactive",
			}) //Empty (No rows)
		mock.ExpectQuery(regexp.QuoteMeta("select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(InboundFileKey as char(36)) as InboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from InboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey  where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1")).WillReturnRows(rows)

		//Execute funcion to be tested
		output := fitted_filesController.GetInboundFiles(mock_ctx)

		//Output should be untampered
		if reflect.DeepEqual(output, []models.Inboundfile{}) {
			t.Fatalf(`Output not expected`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})
}

func TestGetOutboundFiles(t *testing.T) {
	mock, fitted_filesController, err := setupfilesController()
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)

	t.Run("validFileId", func(t *testing.T) {
		sampleId := "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3" //Doesn't really matter because Query output is hardcoded

		//Preparing a Request Body
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("Id", sampleId))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/api/inboundfiles", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Sample DB output, can be anything
		sample_row1 := models.Outboundfile{
			OutboundFileKey: "F3F53EB3-91DD-484E-B38F-0063D03343A8",
			NaesbUserKey:    "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3",
			UsKey:           "D6DFAD93-8B80-4119-8B94-6814F1ED75BD",
			UsCommonCode:    "829244552110T",
			ThemKey:         "64E1BE5C-EDF7-4B29-8455-249D0FCA3DA9",
			ThemCommonCode:  "957877905",
			Filename:        "997_957877905_829244552110T_000000014.edi",
			Plaintext:       "ISA*00*          *00*          *01*829244552110T  *01*957877905      *161027*0314*U*00401*000000016*0*T*>~GS*FA*8292445521100*957877905*20161027*0314*16*X*004010~ST*997*129640000~AK1*GE*25~AK2*814*0002~AK5*A~AK9*A*1*1*1~SE*6*129640000~GE*1*16~IEA*1*000000016~",
			Ciphertext: sql.NullString{
				String: "-----BEGIN PGP MESSAGE-----\r\nVersion: GnuPG v1.4.7 (MingW32)\r\n\r\nhIwD4t+GNWa8QGsBA/0eNKwgz+lfDYiXRDwhqqn9DHWDf2uvG0PFkL1fyO53t6te\r\ntcKR04Maow2IPgJK6lQgXGgCtyNMtfdt5z041ahawF7DTaiUrxIikgg4G6IESmD9\r\nmpnYyjP/aRk1wtbenY7JGnqV+7ghzzyiLm46PORCwEXf6SHYQZ+in/dywes3m8nA\r\nQE0jypXZZBGN8f2sUn37h6bHSmKQsr/MoSpXYNh6PTJuft+BhigD7ULAXy1PqbG8\r\nCptkoMlE3y/5mFQFeLsUHg3QdGd6vqJ3mrPhYN77hJt634lAsks5FinuKQIbN+Jh\r\nnBKcXhVwvcMfk6zoz3e0BpqUXXZ1bdTj0SNlpOhxVzl5c2ZKNDOCHBHxH4QHfi1F\r\nHdwnzDmGnoJnwo39sTGhvFFp6dG7vNXJ5MauGQmZyFnAIVnwdgA7/2Qrymm2P0Fe\r\nYXkQc5cZ0yf7vsmV0yQ9aERuI9FAy6mBgwCcvLUiRzqB2zMLrjhpap0+RA3eANF4\r\nqFYKRtqgcEeM4WQjYKV8TaQ=\r\n=RsD7\r\n-----END PGP MESSAGE-----\r\n",
				Valid:  true,
			},
			Attempt1At: sql.NullString{
				String: "2016-10-27T15:17:52.313Z",
				Valid:  true,
			},
			Attempt2At: sql.NullString{
				String: "",
				Valid:  false,
			},
			Attempt3At: sql.NullString{
				String: "",
				Valid:  false,
			},
			Receipt: sql.NullString{
				String: "730724610071040",
				Valid:  true,
			},
			Result: sql.NullString{
				String: "--GISB7866\r\nContent-type: text/html\r\n\r\n<HTML><HEAD><TITLE>Acknowledgement Receipt</TITLE></HEAD><BODY><P>\r\ntime-c=20161027151752*\r\nrequest-status=ok*\r\nserver-id=emkt-naesb.centerpointenergy.com*\r\ntrans-id=730724610071040*\r\n</P></BODY></HTML>\r\n--GISB7866\r\nContent-type: text/plain\r\n\r\ntime-c=20161027151752*\r\nrequest-status=ok*\r\nserver-id=emkt-naesb.centerpointenergy.com*\r\ntrans-id=730724610071040*\r\n--GISB7866--",
				Valid:  true,
			},
			Escalated: "false",
			EscalatedAt: sql.NullString{
				String: "",
				Valid:  false,
			},
			Debug:     "false",
			CreatedAt: "2016-10-27T15:14:45.173Z",
			EmpowerOutgoingEdiFileKey: sql.NullString{
				String: "��2��FiK���^�'��",
				Valid:  true,
			},
			DoNotSend: "false",
			LastLocation: sql.NullString{
				String: "H:\\FlightTest\\NAESB\\_Files\\FT_Illuminar_Alltex\\OUT",
				Valid:  true,
			},
			Ciphered:       "true",
			Posted:         "true",
			OutboundFileId: "16",
			Inactive:       "false",
		}
		sample_row2 := models.Outboundfile{
			OutboundFileKey: "A0D6BE4E-39AF-4FDE-977F-087910BD9137",
			NaesbUserKey:    "8B0528AB-6E22-40E2-9B60-A4A6C584E6E3",
			UsKey:           "D6DFAD93-8B80-4119-8B94-6814F1ED75BD",
			UsCommonCode:    "829244552110T",
			ThemKey:         "D5673342-0082-4281-875D-D92B17143CB3",
			ThemCommonCode:  "103994067",
			Filename:        "997_103994067_829244552110T_000137059.edi",
			Plaintext:       "ISA*00*          *00*          *01*829244552110T  *01*103994067      *161027*0314*U*00401*000000012*0*T*>~GS*FA*8292445521100*1039940674000*20161027*0314*12*X*004010~ST*997*837370000~AK1*GE*1539~AK2*814*0003~AK5*A~AK9*A*1*1*1~SE*6*837370000~GE*1*12~IEA*1*000000012~",
			Ciphertext: sql.NullString{
				String: "-----BEGIN PGP MESSAGE-----\r\nVersion: GnuPG v1.4.7 (MingW32)\r\n\r\nhIwDKpwAYzBCHt8BBACebjXV2M+XOm07tAAq3fU821oUIzR5/DK9OpQ8m/mSJv8a\r\nboqApaju7xq6XyPsEEdgrEzbJuym+5USv0JTQluiyjvRQ80bptyvc+/s7Ro5YpSV\r\nL2PJtENI7Gc2a+qx07BuQgaQ5MZnByGnCFztqwRnORZ5VCnoeSerJdQpzg2uPNLA\r\nXAGdU2MOutlhgonZN0W7g75xX1HAs2u50biyNQ96GziR04ugCrCoZScARLaBODaw\r\ngaMEiihDxIh+JdlYkLgQrf9f+lj/B3k/uPpAbCO9r0VWfcQLYnc3/gNL9gpT1CLo\r\n7adSfnvgdVtHSzBPxzzDYa80z44xmLr3YRrcQ5eV45qiaNiDUsWaYW/IneHvIhCc\r\nsT++353AMyvOk+4TnttAlIWa2QMNzdOwKtilouZl+oEPet73WvcbDvtc8lmBWkjJ\r\n7pdoFe1G17M9PaZ8xU507LDsYshymT7GKhBEEEZ2kQ89V+k8/9l82q7bX3Lw6c7O\r\nLvlYWzXPEOc5+lDQyhyzoDlKZEPx4xvd0Wi3c7iIKgT+JmC/aEa/ip+9WBSE\r\n=FV9k\r\n-----END PGP MESSAGE-----\r\n",
				Valid:  true,
			},
			Attempt1At: sql.NullString{
				String: "2016-10-27T15:17:49.31Z",
				Valid:  true,
			},
			Attempt2At: sql.NullString{
				String: "",
				Valid:  false,
			},
			Attempt3At: sql.NullString{
				String: "",
				Valid:  false,
			},
			Receipt: sql.NullString{
				String: "477599469120064",
				Valid:  true,
			},
			Result: sql.NullString{
				String: "------=_NAESB_Report_edibzlmkt01.test.corp.oncor.com-12-1474997419726-3.\r\nContent-type: text/html\r\n\r\n<HTML><HEAD><TITLE>Acknowledgement Receipt Success</TITLE></HEAD><BODY><PRE>\r\ntime-c=20161027151750*\r\nrequest-status=ok*\r\nserver-id=edibzlmkt01.test.corp.oncor.com*\r\ntrans-id=477599469120064*\r\n</PRE></BODY></HTML>\r\n------=_NAESB_Report_edibzlmkt01.test.corp.oncor.com-12-1474997419726-3.\r\nContent-type: text/plain\r\n\r\ntime-c=20161027151750*\r\nrequest-status=ok*\r\nserver-id=edibzlmkt01.test.corp.oncor.com*\r\ntrans-id=477599469120064*\r\n------=_NAESB_Report_edibzlmkt01.test.corp.oncor.com-12-1474997419726-3.--\r\n",
				Valid:  true,
			},
			Escalated: "false",
			EscalatedAt: sql.NullString{
				String: "",
				Valid:  false,
			},
			Debug:     "false",
			CreatedAt: "2016-10-27T15:14:43.887Z",
			EmpowerOutgoingEdiFileKey: sql.NullString{
				String: "����/�F�y��S�\f",
				Valid: true,
			},
			DoNotSend: "false",
			LastLocation: sql.NullString{
				String: "H:\\FlightTest\\NAESB\\_Files\\FT_Illuminar_Alltex\\OUT",
				Valid:  true,
			},
			Ciphered:       "true",
			Posted:         "true",
			OutboundFileId: "206",
			Inactive:       "false",
		}

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"OutboundFileKey",
				"NaesbUserKey",
				"UsKey",
				"UsCommonCode",
				"ThemKey",
				"ThemCommonCode",
				"Filename",
				"Plaintext",
				"Ciphertext",
				"Attempt1At",
				"Attempt2At",
				"Attempt3At",
				"Receipt",
				"Result",
				"Escalated",
				"EscalatedAt",
				"Debug",
				"CreatedAt",
				"EmpowerOutgoingEdiFileKey",
				"DoNotSend",
				"LastLocation",
				"Ciphered",
				"Posted",
				"OutboundFileId",
				"Inactive",
			}).
			AddRow(testutils.RowFromStruct(sample_row1)...).
			AddRow(testutils.RowFromStruct(sample_row2)...)
		mock.ExpectQuery(regexp.QuoteMeta("select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(OutboundFileKey as char(36)) as OutboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from OutboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1")).WillReturnRows(rows)

		//Execute funcion to be tested
		output := fitted_filesController.GetOutboundFiles(mock_ctx)

		//Output should be untampered
		if !reflect.DeepEqual(output, []models.Outboundfile{sample_row1, sample_row2}) {
			t.Fatalf(`Output doesn't match`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})

	t.Run("invalidFileId", func(t *testing.T) {
		sampleId := "69420" //Doesn't really matter because Query output is hardcoded

		//Preparing a Request Body
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		assert.NoError(t, mw.WriteField("Id", sampleId))
		w, err := mw.CreateFormFile("file", "test")
		if assert.NoError(t, err) {
			_, err = w.Write([]byte("test"))
			assert.NoError(t, err)
		}
		mw.Close()

		//Preparing a Context containing Request Body
		mock_ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		mock_ctx.Request = httptest.NewRequest("POST", "http://localhost:8000/api/inboundfiles", buf)
		mock_ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

		//Mocking DB response
		rows := sqlmock.NewRows(
			[]string{
				"OutboundFileKey",
				"NaesbUserKey",
				"UsKey",
				"UsCommonCode",
				"ThemKey",
				"ThemCommonCode",
				"Filename",
				"Plaintext",
				"Ciphertext",
				"Attempt1At",
				"Attempt2At",
				"Attempt3At",
				"Receipt",
				"Result",
				"Escalated",
				"EscalatedAt",
				"Debug",
				"CreatedAt",
				"EmpowerOutgoingEdiFileKey",
				"DoNotSend",
				"LastLocation",
				"Ciphered",
				"Posted",
				"OutboundFileId",
				"Inactive",
			}) //Empty (No rows)
		mock.ExpectQuery(regexp.QuoteMeta("select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(OutboundFileKey as char(36)) as OutboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from OutboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1")).WillReturnRows(rows)

		//Execute funcion to be tested
		output := fitted_filesController.GetOutboundFiles(mock_ctx)

		//Output should be untampered
		if reflect.DeepEqual(output, []models.Outboundfile{}) {
			t.Fatalf(`Output not expected`)
		}

		//Output should come from a DB query only.
		if eror := mock.ExpectationsWereMet(); eror != nil {
			t.Fatalf(eror.Error())
		}
	})
}
