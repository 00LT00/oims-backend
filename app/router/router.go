package router

import (
	"github.com/gin-gonic/gin"
	"oims/service"
	"oims/utils"
)

var s = service.Service

func InitRouter() {
	r := s.Engine
	r.Use(check, recovery)
	r.POST("/image", f(getJpeg))

	r.GET("/xml", f(getXml))

	_ = r.Run(conf.Server.Port)
}

func f(h func(c *gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := h(c)
		c.JSON(utils.MakeSuccessJSON(data))
	}
}
