package handles

import (
	"MyLNPU/internal/logger"
	"MyLNPU/internal/service"
	"MyLNPU/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetExpTable(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	table, err := service.GetExpTable(openid)
	if err != nil {
		logger.Errorf("获取实验课程信息失败... %s", err)
		return nil, err
	}
	logger.Println("获取实验课程信息成功 [%s]", openid)
	return map[string]any{"table": table}, nil

}
