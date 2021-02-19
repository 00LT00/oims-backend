package service

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type service struct {
	DB     *gorm.DB
	Engine *gin.Engine
}
