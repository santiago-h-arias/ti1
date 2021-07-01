package testutils

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"

	models "tinc1/Models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

//Converts a Struct into []driver.Value so that
// it can be passed into (*sqlmock.Rows).AddRow(values ...driver.Value)
func RowFromStruct(sample interface{}) []driver.Value {
	var ret []driver.Value

	rv := reflect.ValueOf(sample)
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		dv := driver.Value(fv.Interface())
		ret = append(ret, dv)
	}

	return ret
}

//For Mocking JWTService Interface
type mock_jwtService struct {
	secretkey string
	issuer    string
}

func (jwtService *mock_jwtService) GenerateToken(email string) string {
	if email == "ajith@thinkbridge.in" {
		return "1234_valid_token_4321"
	}

	return ""
}

func (jwtService *mock_jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return &jwt.Token{}, nil
}

func NewMock_JWTService() *mock_jwtService {
	return &mock_jwtService{
		secretkey: "UjgFm344XW",
		issuer:    "thinkbridgeIdProvider",
	}
}

//For Mocking Dao Interface (for LoginService)
type mock_dao struct {
	db sqlx.DB
}

func NewMock_Dao(db sqlx.DB) *mock_dao {
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

//Implementing Try Catch in Golang
/*

type Exception interface{}

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

func Throw(e Exception) {
	panic(e)
}

func (b Block) Do() {
	if b.Finally != nil {
		defer b.Finally()
	}

	if b.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				b.Catch(r)
			}
		}()
	}

	b.Try()
}

func TestNewDao(t *testing.T) {
	Block{
		Try: func() {
			_ = NewDao()
		},
		Catch: func(e Exception) {
			fmt.Printf("Caught Error : %v\n", e)
		},
	}.Do()
}
*/
