package common

import (
	"MyLNPU/internal/logger"
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
	logger.Errorf(err.Error())
	c.Abort()
}

func ErrorStrResp(c *gin.Context, code int, str string) {
	c.JSON(200, Resp[any]{
		Code:    code,
		Message: str,
		Data:    nil,
	})
	logger.Errorf(str+" Path: %s IP: %s", c.Request.URL, c.ClientIP())
	c.Abort()
}
