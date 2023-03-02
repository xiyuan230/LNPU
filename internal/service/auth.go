package service

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/constant"
	"MyLNPU/internal/db"
	"MyLNPU/internal/log"
	"MyLNPU/internal/model"
	"MyLNPU/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

func Login(code string) (string, error) {
	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", constant.APPID, constant.APPSECRET, code)
	rep, err := http.Get(wxUrl)
	if err != nil {
		log.Errorf("获取openid失败... %s", err)
		return "", err
	}
	defer rep.Body.Close()
	result, err := io.ReadAll(rep.Body)
	if err != nil {
		log.Errorf("响应结果解析失败... %s", err)
		return "", err
	}
	wxResult := model.WXLoginRequest{}
	json.Unmarshal(result, &wxResult)
	_, err = db.GetUserByID(wxResult.Openid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u := model.User{}
			u.OpenID = wxResult.Openid
			err := db.CreateUser(&u)
			if err != nil {
				log.Errorf("插入User失败... %s", err)
				return "", err
			}
		} else {
			log.Errorf("查询User失败... %s", err)
			return "", err
		}
	}
	token, err := utils.JWTNewToken(wxResult.Openid)
	if err != nil {
		return "", err
	}
	cache.Set("lnpu:token:"+token, wxResult.Openid, time.Hour*2)
	return token, nil
}

func CheckTokenExpiration(token string) bool {
	isExpiration := utils.CheckTokenStatus(token)
	return isExpiration
}
