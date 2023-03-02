package model

type BaseUser struct {
	UserName     string `json:"userName"`
	UserPassword string `json:"userPassword"`
}

type WXLoginRequest struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
}
