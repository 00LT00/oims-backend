package service

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	_ "oims/service/gpu"
	"os"
	"path/filepath"
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
	Conf.Path.Result = dir + Conf.Path.Result
	Conf.Path.History = dir + Conf.Path.History
	Conf.Path.Log = dir + Conf.Path.Log
	initPath()
	initLog()
	r := initGin()
	Service.Engine = r

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

func initPath() {
	err := os.MkdirAll(Conf.Path.History, os.ModePerm)
	if err != nil {
		panic(err)
	}
	//os.Chmod(Conf.Path.History, os.ModePerm)
	err = os.MkdirAll(Conf.Path.Result, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(Conf.Path.Log, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func initLog() {
	logFile, err := os.OpenFile(filepath.Join(Conf.Path.Log, "log.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	errlogFile, err := os.OpenFile(filepath.Join(Conf.Path.Log, "errlog.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Service.Logger = log.New(logFile, "", log.LstdFlags)
	Service.ErrLogger = log.New(errlogFile, "", log.LstdFlags)
	Service.Logger.Println("init")
	Service.ErrLogger.Println("init")
}
