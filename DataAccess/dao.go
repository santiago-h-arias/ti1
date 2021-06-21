package dataaccess

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	config "tinc1/Config"
	models "tinc1/Models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dao interface {
	CheckUser(userId string, password string) (bool, models.NaesbUser)
}

type Files interface {
	Get_files(id string) []models.Inboundfile
}
type dao struct {
	db *gorm.DB
}

func NewDao() Dao {
	configuration, conferr := GetConfiguration()
	if conferr != nil {
		log.Fatalln(conferr)
		return nil
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		configuration.User,
		configuration.Password,
		configuration.DBAddress,
		configuration.Port,
		configuration.DBName,
	)

	db, dberr := gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if dberr != nil {
		log.Fatalln(dberr)
		return nil
	}

	return &dao{
		db: db,
	}
}
func NewFiles() Files {
	configuration, conferr := GetConfiguration()
	if conferr != nil {
		log.Fatalln(conferr)
		return nil
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		configuration.User,
		configuration.Password,
		configuration.DBAddress,
		configuration.Port,
		configuration.DBName,
	)

	db, dberr := gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
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
	result := dao.db.Table("NaesbUser").Where("Email = ? AND Password = ?", email, password).Find(&user)
	return result.RowsAffected > 0, user
}

func (dao *dao) Get_files(id string) []models.Inboundfile {
	var findfiles []models.Inboundfile
	//data := dao.db.Table("InboundFiles").Select("*").Joins("join NaesbUserUs on NaesbUserUs.Uskey = InboundFiles.Uskey").Where("NaesbUserUs.Inactive= ?", "0").Where("NaesbUserUs.NaesbUserKey = ?", id).Find(&findfiles)
	//dao.db.Table("InboundFiles").Where("UsKey = ?", id).Find(&findfiles)
	data := dao.db.Raw("select* from InboundFiles where UsKey='D6DFAD93-8B80-4119-8B94-6814F1ED75BD'").Scan(&findfiles)
	fmt.Println(data)
	return findfiles

}
