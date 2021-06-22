package dataaccess

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	config "tinc1/Config"
	models "tinc1/Models"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type Dao interface {
	CheckUser(userId string, password string) (bool, models.NaesbUser)
}

type dao struct {
	db *sqlx.DB
}

func NewDao() Dao {
	configuration, conferr := GetConfiguration()
	if conferr != nil {
		log.Fatalln(conferr)
		return nil
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s/?database=%s",
		configuration.User,
		configuration.Password,
		configuration.DBAddress,
		configuration.Port,
		configuration.DBName,
	)
	db, dberr := sqlx.Connect("sqlserver", dsn)
	if dberr != nil {
		log.Fatalln(dberr)
		return nil
	}

	return &dao{
		db: db,
	}
}

func GetConfiguration() (config.DBConfig, error) {
	config := config.DBConfig{}
	file, err := os.Open("./configuration.json")
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (dao *dao) CheckUser(email string, password string) (bool, models.NaesbUser) {
	var user models.NaesbUser

	err := dao.db.QueryRowx("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2", email, password).StructScan(&user)
	if err == nil {
		return true, user
	}
	return false, user

}
