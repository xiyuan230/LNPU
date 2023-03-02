package common

import (
	"MyLNPU/internal/log"
	"github.com/gin-gonic/gin"
)

func SuccessResp(c *gin.Context, data ...any) {
	if len(data) == 0 {
		c.JSON(200, Resp[any]{
			Code:    200,
			Message: "success",
			Data:    nil,
		})
		return
	}
	c.JSON(200, Resp[any]{
		Code:    200,
		Message: "success",
		Data:    data[0],
	})
}

func ErrorResp(c *gin.Context, code int, err error) {
	c.JSON(200, Resp[any]{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	})
	log.Errorf(err.Error())
	c.Abort()
}

func ErrorStrResp(c *gin.Context, code int, str string) {
	c.JSON(200, Resp[any]{
		Code:    code,
		Message: str,
		Data:    nil,
	})
	log.Errorf(str+" Path: %s IP: %s", c.FullPath(), c.ClientIP())
	c.Abort()
}
