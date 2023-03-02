package handles

import (
	"MyLNPU/api/common"
	"MyLNPU/internal/cache"
	"MyLNPU/internal/service"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	value := c.Query("code")
	if value == "" {
		common.ErrorStrResp(c, 401, "非法请求")
		return
	}
	token, err := service.Login(value)
	if err != nil {
		common.ErrorResp(c, 500, err)
		return
	}

	common.SuccessResp(c, map[string]any{"token": token})
}

func CheckTokenStatus(c *gin.Context) {
	token := c.GetHeader("Authorization")
	status := service.CheckTokenExpiration(token)
	if !status {
		cache.Del("lnpu:token:" + token)
	}
	common.SuccessResp(c, map[string]any{"status": status})
}
