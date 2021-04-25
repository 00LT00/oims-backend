package service

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type service struct {
	DB        *gorm.DB
	Engine    *gin.Engine
	Logger    *log.Logger
	ErrLogger *log.Logger
}
