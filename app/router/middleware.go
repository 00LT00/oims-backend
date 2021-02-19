package router

import (
	"github.com/gin-gonic/gin"
	"oims/error"
	"oims/service"
	"oims/utils"
)

var conf = service.Conf

func check(c *gin.Context) {
	sign := c.GetHeader("sign")
	if sign != conf.Server.Sign {
		c.AbortWithStatusJSON(utils.MakeErrJSON(403, "40300", "forbidden"))
	}
}

func recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error.HttpError); ok {
				c.AbortWithStatusJSON(utils.MakeErrJSON(e.HttpStatusCode, e.Err.ErrCode, e.Err.Msg))
				return
			}
			c.AbortWithStatusJSON(utils.MakeErrJSON(500, "50000", r))
			return
		}
	}()
	c.Next()
}
