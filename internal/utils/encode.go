package utils

import (
	"encoding/base64"
	"strconv"
	"strings"
)

func EncodeByBase64(userAccount, userPassword string) string {
	username := base64.StdEncoding.EncodeToString([]byte(userAccount))
	password := base64.StdEncoding.EncodeToString([]byte(userPassword))
	encode := username + "%%%" + password
	return encode
}

func EncodeByUrlParam(userAccount, userPassword, str string) string {
	split := strings.Split(str, "#")
	sCode := split[0]
	sxh := split[1]
	code := userAccount + "%%%" + userPassword
	encoded := ""
	for i := 0; i < len(code); i++ {
		if i < 20 {
			t, _ := strconv.Atoi(sxh[i : i+1])
			encoded += code[i:i+1] + sCode[0:t]
			sCode = sCode[t:len(sCode)]
		} else {
			encoded += code[i:len(code)]
			i = len(code)
		}
	}
	return encoded
}
