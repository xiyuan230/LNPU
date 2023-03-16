package service

import (
	"MyLNPU/internal/errs"
	"MyLNPU/internal/logger"
	"MyLNPU/internal/utils"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func loginWithSSO(userName, password string) (*http.Client, error) {
	client, err := utils.NewHttpClient()
	if err != nil {
		logger.Errorf("创建HttpClient时出错... %s", err)
		return nil, err
	}
	ssoResp, err := client.Get(SSOLoginURl)
	if err != nil {
		return nil, err
	}
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

func loginWithJwxt(userAccount, userPassword string) (*http.Client, error) {
	for i := 0; i < 3; i++ {
		code, client, err := utils.GetVerifyCode("https://jwxt.lnpu.edu.cn/jsxsd/verifycode.servlet")
		if err != nil {
			if errors.Is(err, image.ErrFormat) {
				continue
			}
			return nil, err
		}
		encode := utils.EncodeByBase64(userAccount, userPassword)
		values := url.Values{}
		values.Add("RANDOMCODE", code)
		values.Add("encoded", encode)
		values.Add("userAccount", "")
		payload := values.Encode()
		request, _ := http.NewRequest("POST", "https://jwxt.lnpu.edu.cn/jsxsd/xk/LoginToXk", strings.NewReader(payload))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response, err := client.Do(request)
		if err != nil {
			logger.Println("教务系统登录失败...")
			continue
		}
		body, _ := io.ReadAll(response.Body)
		if strings.Contains(string(body), "密码错误") {
			return nil, errs.ErrPasswordWrong
		} else if strings.Contains(string(body), "验证码错误") {
			logger.Println("验证码错误...")
			continue
		}
		//cookies := response.Request.Header.Values("Cookie")
		//cookie := strings.Join(cookies, "")
		return client, nil
	}
	return nil, errs.ErrJwxtLoginFailed
}
