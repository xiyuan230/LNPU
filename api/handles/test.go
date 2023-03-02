package handles

import (
	"MyLNPU/internal/log"
	"MyLNPU/internal/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func TestSSO(c *gin.Context) {
	const USERNAME = "2012200108"
	const PASSWORD = "Sql123.."
	client, err := utils.NewHttpClient()
	if err != nil {
		log.Errorf("%s", err)
	}
	resp, _ := client.Get("https://sso.lnpu.edu.cn/login")
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	croypt := doc.Find("#login-croypto").First().Text()
	flowkey := doc.Find("#login-page-flowkey").First().Text()
	values := url.Values{}
	values.Add("username", "2012200108")
	values.Add("type", "UsernamePassword")
	values.Add("_eventId", "submit")
	values.Add("geolocation", "")
	values.Add("execution", flowkey)
	values.Add("croypto", croypt)
	values.Add("password", utils.DesECBEncrypt("Sql123..", croypt))
	encode := values.Encode()
	req, _ := http.NewRequest("POST", "https://sso.lnpu.edu.cn/login", strings.NewReader(encode))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "sso.lnpu.edu.cn")
	req.Header.Add("Origin", "https://sso.lnpu.edu.cn")
	req.Header.Add("Referer", "https://sso.lnpu.edu.cn/login?service=https:%2F%2Fjwxt.lnpu.edu.cn%2Fjsxsd%2Fsso.jsp")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	rep, _ := client.Do(req)
	all, _ := io.ReadAll(rep.Body)
	defer rep.Body.Close()
	get, _ := client.Get("https://sso.lnpu.edu.cn/login?service=https%3A%2F%2Fjwxt.lnpu.edu.cn%2Fjsxsd%2Fsso.jsp")
	defer get.Body.Close()
	readAll, _ := io.ReadAll(get.Body)
	fmt.Println(string(readAll))
	fmt.Println(get.Request.Header.Values("Cookie")[0])
	c.String(200, string(all))
}
