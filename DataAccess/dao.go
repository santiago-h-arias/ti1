package dataaccess

import (
	"fmt"
	"log"

	models "tinc1/Models"
	dbutils "tinc1/Utils"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type Dao interface {
	CheckUser(userId string, password string) (bool, models.NaesbUser)
	GetInboundFiles(id string) []models.Inboundfile
}

type dao struct {
	db *sqlx.DB
}

func NewDao() Dao {
	configuration, conferr := dbutils.GetConfiguration()
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

func (dao *dao) CheckUser(email string, password string) (bool, models.NaesbUser) {
	var user models.NaesbUser

	err := dao.db.QueryRowx("select cast(NaesbUserKey as char(36)) as NaesbUserKey, Name, Email from NaesbUser where Email=@p1 and Password=@p2", email, password).StructScan(&user)
	if err == nil {
		return true, user
	}
	return false, user

}

func (dao *dao) GetInboundFiles(id string) []models.Inboundfile {
	var findfiles []models.Inboundfile
	dao.db.Select(&findfiles, "select *, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(InboundFileKey as char(36)) as InboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(ThemKey as char(36)) as ThemKey from InboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey  where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1", id)
	return findfiles

}
