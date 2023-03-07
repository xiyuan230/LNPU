package utils

import (
	"MyLNPU/internal/logger"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var signKey = []byte("Sql123..")

type MyClaims struct {
	OpenID string `json:"openid"`
	jwt.RegisteredClaims
}

func JWTNewToken(openid string) (string, error) {
	claims := MyClaims{
		OpenID: openid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			Issuer:    "xiyuan",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)
	if err != nil {
		logger.Errorf("Token签发失败... %s", err)
		return "", err
	}
	return token, nil
}

func JWTParseToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}
	claims := MyClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("unexpected signing method")
		}
		return signKey, nil
	})
	if err != nil {
		return "", err
	}
	return claims.OpenID, nil
}

func CheckTokenStatus(token string) bool {
	claims := MyClaims{}
	jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("unexpected signing method")
		}
		return signKey, nil
	})
	return claims.ExpiresAt.Unix() > time.Now().Unix()
}
