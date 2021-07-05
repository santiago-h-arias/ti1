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
	GetOutboundFiles(id string) []models.Outboundfile
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
	err := dao.db.Select(&findfiles, "select if2.UsCommonCode, if2.ThemCommonCode , if2.Filename, if2.Plaintext, if2.Ciphertext, if2.ReceivedAt, if2.TransactionId, if2.Processed, if2.InboundFileId, u.Name as UsName, t.Name  as ThemName, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(InboundFileKey as char(36)) as InboundFileKey, cast(if2.UsKey as char(36)) as UsKey, cast(if2.ThemKey as char(36)) as ThemKey from InboundFiles if2 left join NaesbUserUs nuu on nuu.UsKey = if2.Uskey  join Us u on u.UsKey = if2.UsKey join Them t on t.ThemKey = if2.ThemKey where nuu.Inactive = 0 and nuu.NaesbUserKey=@p1", id)
	if err != nil {
		fmt.Println(err)
	}
	return findfiles
}

func (dao *dao) GetOutboundFiles(id string) []models.Outboundfile {
	var findfiles []models.Outboundfile
	err := dao.db.Select(&findfiles, "select of2.UsCommonCode, of2.ThemCommonCode, of2.Filename, of2.Plaintext, of2.Ciphertext, of2.Attempt1At, of2.Attempt2At, of2.Attempt3At,of2.Receipt, of2.Result, of2.Escalated , of2.EscalatedAt , of2.Debug , of2.CreatedAt , of2.EmpowerOutgoingEdiFileKey , of2.DoNotSend , of2.LastLocation , of2.Ciphered , of2.Posted, u.Name as UsName, t.Name  as ThemName, cast(nuu.NaesbUserKey as char(36)) as NaesbUserKey, cast(OutboundFileKey as char(36)) as OutboundFileKey, cast(of2.UsKey as char(36)) as UsKey, cast(of2.ThemKey as char(36)) as ThemKey from OutboundFiles of2 left join NaesbUserUs nuu on nuu.UsKey = of2.Uskey join Us u on u.UsKey = of2.UsKey join Them t on t.ThemKey = of2.ThemKey where nuu.Inactive = 0 and nuu.NaesbUserKey =@p1", id)
	if err != nil {
		fmt.Println(err)
	}
	return findfiles
}
