package handles

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/errs"
	"MyLNPU/internal/service"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) (any, error) {
	value := c.Query("code")
	if value == "" {
		return nil, errs.ErrParamMiss
	}
	token, err := service.Login(value)
	if err != nil {
		return nil, err
	}
	return map[string]any{"token": token}, nil
}

func CheckTokenStatus(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	status := service.CheckTokenExpiration(token)
	if !status {
		cache.Del("lnpu:token:" + token)
	}
	return map[string]any{"status": status}, nil
}
