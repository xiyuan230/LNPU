package handles

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/errs"
	"MyLNPU/internal/model"
	"MyLNPU/internal/service"
	"MyLNPU/internal/utils"
	"github.com/gin-gonic/gin"
	"time"
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

func BindUser(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = c.ShouldBind(&user)
	if err != nil {
		return nil, err
	}
	pattern := c.Param("pattern")
	switch pattern {
	case "jwxt":
		if user.StudentID == "" || user.JwxtPassword == "" {
			return "", errs.ErrUserIllegal
		}
		cookie, err := service.JwxtLoginBindJwxt(user.StudentID, user.JwxtPassword)
		if err != nil {
			return nil, err
		}
		u := model.User{OpenID: openid, StudentID: user.StudentID, JwxtPassword: user.JwxtPassword}
		err = service.UpdateUser(&u)
		if err != nil {
			return nil, err
		}
		cache.Set("lnpu:jwxt:cookie:"+openid, cookie, time.Hour*1)
		stu, err := service.GetStudentInfo(openid)
		if err != nil {
			return nil, err
		}
		return map[string]any{"student_info": stu}, nil
	case "sso":
		if user.StudentID == "" || user.SSOPassword == "" {
			return "", errs.ErrUserIllegal
		}
		cookie, err := service.JwxtLoginBindSSO(user.StudentID, user.SSOPassword)
		if err != nil {
			return nil, err
		}
		u := model.User{OpenID: openid, StudentID: user.StudentID, SSOPassword: user.SSOPassword}
		err = service.UpdateUser(&u)
		if err != nil {
			return nil, err
		}
		cache.Set("lnpu:jwxt:cookie:"+openid, cookie, time.Hour*1)
		stu, err := service.GetStudentInfo(openid)
		if err != nil {
			return nil, err
		}
		return map[string]any{"student_info": stu}, nil
	case "experiment":
		if user.StudentID == "" || user.ExpPassword == "" {
			return "", errs.ErrUserIllegal
		}
		cookie, err := service.ExpLoginBind(user.StudentID, user.ExpPassword)
		if err != nil {
			return nil, err
		}
		u := model.User{OpenID: openid, StudentID: user.StudentID, ExpPassword: user.ExpPassword}
		err = service.UpdateUser(&u)
		if err != nil {
			return nil, err
		}
		cache.Set("lnpu:exp:cookie:"+openid, cookie, time.Hour*2)
		return nil, nil
	default:
		return nil, errs.ErrPathIllegal
	}
	return nil, errs.ErrPathIllegal
}

func CheckTokenStatus(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	status := service.CheckTokenExpiration(token)
	if !status {
		cache.Del("lnpu:token:" + token)
	}
	return map[string]any{"status": status}, nil
}
