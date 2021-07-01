package testutils

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

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
	secretKey string
	issuer    string
}

type mock_jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func (jwtService *mock_jwtService) GenerateToken(email string) string {
	// Set custom and standard claims
	claims := &mock_jwtCustomClaims{
		email,
		jwt.StandardClaims{
			Issuer:    jwtService.issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(jwtService.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtService *mock_jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(jwtService.secretKey), nil
	})
}

func NewMock_JWTService() *mock_jwtService {
	return &mock_jwtService{
		secretKey: "UjgFm344XW",
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
