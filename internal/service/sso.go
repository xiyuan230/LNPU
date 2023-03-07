package service

import (
	"MyLNPU/internal/errs"
	"MyLNPU/internal/logger"
	"MyLNPU/internal/utils"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

func SSOLogin(userName, password string) (*http.Client, error) {
	client, err := utils.NewHttpClient()
	if err != nil {
		logger.Errorf("创建HttpClient时出错... %s", err)
		return nil, err
	}
	ssoResp, _ := client.Get(SSOLoginURl)
	defer ssoResp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(ssoResp.Body)
	croypt := doc.Find("#login-croypto").First().Text()
	flowkey := doc.Find("#login-page-flowkey").First().Text()
	values := url.Values{}
	values.Add("username", userName)
	values.Add("type", "UsernamePassword")
	values.Add("_eventId", "submit")
	values.Add("geolocation", "")
	values.Add("execution", flowkey)
	values.Add("croypto", croypt)
	values.Add("password", utils.DesECBEncrypt(password, croypt))
	encode := values.Encode()
	ssoLoginReq, _ := http.NewRequest("POST", SSOLoginURl, strings.NewReader(encode))
	ssoLoginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ssoLoginReq.Header.Add("Host", "sso.lnpu.edu.cn")
	ssoLoginReq.Header.Add("Origin", "https://sso.lnpu.edu.cn")
	ssoLoginReq.Header.Add("Referer", JwxtLoginUrlWithSSO)
	ssoLoginReq.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	ssoLoginResp, err := client.Do(ssoLoginReq)
	if err != nil {
		logger.Errorf("统一认证登录请求失败... %s", err)
		return nil, err
	}
	defer ssoLoginResp.Body.Close()
	if ssoLoginResp.Request.URL.String() == "https://sso.lnpu.edu.cn/login" {
		return nil, errs.ErrPasswordWrong
	}
	return client, nil
}
