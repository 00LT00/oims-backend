package service

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Service *service
var Conf *conf

func init() {
	Service = new(service)
	dir, _ := os.Getwd()
	_, err := toml.DecodeFile(dir+"/"+"config/config.toml", &Conf)
	if err != nil {
		panic(err)
	}
	//db, err := initDB(Conf)
	//if err != nil {
	//	panic(err)
	//}
	//Service.DB = db
	r := initGin()
	Service.Engine = r
	err = os.MkdirAll("./historys", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func initDB(conf *conf) (*gorm.DB, error) {
	var dsn string
	if os.Getenv("mode") == "Zer0kiriN" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", conf.DB.User, conf.DB.Pass, conf.DB.Addr, conf.DB.DBName)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", conf.DBDev.User, conf.DBDev.Pass, conf.DBDev.Addr, conf.DBDev.DBName)
	}

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func initGin() *gin.Engine {
	return gin.Default()
}
