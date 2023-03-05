package handles

import (
	"MyLNPU/internal/log"
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
		log.Errorf("获取实验课程信息失败... %s", err)
		return nil, err
	}
	log.Println("获取实验课程信息成功 [%s]", openid)
	return map[string]any{"table": table}, nil

}
