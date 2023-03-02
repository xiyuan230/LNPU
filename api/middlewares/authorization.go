package middlewares

import (
	"MyLNPU/api/common"
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func AuthorizationWithToken(c *gin.Context) {
	if c.FullPath() == "/auth/login" {
		c.Next()
		return
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		common.ErrorStrResp(c, 401, "非法请求")
		c.Abort()
		return
	}
	_, err := cache.Get("lnpu:token:" + token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			openid, err := utils.JWTParseToken(token)
			if err != nil {
				common.ErrorResp(c, 401, err)
				c.Abort()
				return
			}
			_, err = db.GetUserByID(openid)
			if err != nil {
				common.ErrorResp(c, 500, err)
				c.Abort()
				return
			}
			cache.Set("lnpu:token:"+token, openid, time.Hour*2)
			c.Next()
			return
		}
	}
	c.Next()
}

func AuthorizationWithAdmin(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		common.ErrorStrResp(c, 401, "非法请求")
		c.Abort()
		return
	}
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		common.ErrorResp(c, 401, err)
		c.Abort()
		return
	}
	user, err := db.GetUserByID(openid)
	if err != nil {
		common.ErrorResp(c, 500, err)
		c.Abort()
		return
	}
	if !user.IsAdmin() {
		common.ErrorStrResp(c, 403, "没有权限")
		c.Abort()
		return
	}
	c.Next()
}
