package service

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/errs"
	"MyLNPU/internal/logger"
	"MyLNPU/internal/model"
	"MyLNPU/internal/utils"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ExpLogin(openid string) (string, error) {
	cache.Del("lnpu:exp:cookie:" + openid)
	user, err := db.GetUserByID(openid)
	if err != nil {
		logger.Errorf("获取用户信息失败... %s", err)
		return "", err
	}
	if user.StudentID == "" || user.ExpPassword == "" {
		return "", errs.ErrUserEmpty
	}
	cookie, err := ExpLoginBind(user.StudentID, user.ExpPassword)
	if err != nil {
		return "", err
	}
	cache.Set("lnpu:exp:cookie:"+openid, cookie, time.Hour*2)
	return cookie, nil
}
func ExpLoginBind(userName, password string) (string, error) {
	client, err := utils.NewHttpClient()
	if err != nil {
		logger.Errorf("创建HttpClient失败... %s", err)
		return "", err
	}
	values := url.Values{}
	values.Add("teaId", userName)
	values.Add("teaPwd", password)
	encode := values.Encode()
	req, _ := http.NewRequest("POST", ExpLoginUrl, strings.NewReader(encode))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(all), "succ") {
		cookie := resp.Header.Get("Set-Cookie")
		return cookie, nil
	} else {
		return "", errs.ErrPasswordWrong
	}
}
func GetExpTable(openid string) (*[]model.Experiment, error) {
	data, err := cache.Get("lnpu:exp:table:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			cookie, err := UpdateExpCookie(openid)
			if err != nil {
				logger.Errorf("登录实践教学平台失败... %s", err)
				return nil, err
			}
			client, err := utils.NewHttpClient()
			if err != nil {
				logger.Errorf("创建HttpClient失败... %s", err)
				return nil, err
			}
			page := 1
			var table []model.Experiment
			for i := 0; i < 10; i++ {
				req, _ := http.NewRequest("GET", ExpTableUrl+"?page="+strconv.Itoa(page), nil)
				req.Header.Add("Cookie", cookie)
				resp, err := client.Do(req)
				if err != nil {
					return nil, err
				}
				logger.Println("正在访问实验列表第%d页", page)
				doc, _ := goquery.NewDocumentFromReader(resp.Body)
				if doc.Find(".links").Length() > 0 {
					return nil, errs.ErrCookieExpire
				}
				if strings.Contains(doc.Find("img[src='images/nomsg.png']").Parent().Text(), "目前没有数据") {
					logger.Println("此页为空")
					break
				}
				courseArea := doc.Find(".page_course_area")
				for j := 0; j < courseArea.Length(); j++ {
					experiment := model.Experiment{}
					experiment.Status = courseArea.Eq(j).Find(".page_course_state").Text()
					experiment.Week = utils.ExpTableWeekHandle(courseArea.Eq(j).Find(".bf_time_img").Parent().Contents().Eq(2).Text())
					experiment.Time = strings.Split(courseArea.Eq(j).Find(".bf_time_img").Text(), " ")[1]
					experiment.Teacher = strings.Split(courseArea.Eq(j).Find(".bf_teacher_img").Text(), "：")[1]
					experiment.Address = strings.Split(courseArea.Eq(j).Find(".bf_pos_img").Text(), "：")[1]
					sec := courseArea.Eq(j).Find(".bf_time_img").Parent().Contents().Eq(3).Text()
					sec = sec[1 : len(sec)-1]
					experiment.Section = sec
					timeText := courseArea.Eq(j).Find(".bf_time_img").Text()
					timeText = strings.Split(strings.Split(timeText, "：")[1], " ")[0]
					parse, _ := time.Parse("2006-01-02", timeText)
					experiment.Date = parse.Unix()
					expName := courseArea.Eq(j).Find(".page_course_state").Parent().Contents().Eq(2).Text()
					expName = strings.Split(strings.TrimSpace(expName), " ")[0]
					experiment.Name = expName
					table = append(table, experiment)
				}
				resp.Body.Close()
				page++
			}
			marshal, _ := json.Marshal(table)
			cache.Set("lnpu:exp:table:"+openid, marshal, time.Hour*4)
			return &table, nil
		}
		return nil, err
	}
	var table []model.Experiment
	json.Unmarshal([]byte(data), &table)
	return &table, nil
}

func UpdateExpCookie(openid string) (string, error) {
	cookie, err := cache.Get("lnpu:exp:cookie:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			login, err := ExpLogin(openid)
			if err != nil {
				return "", err
			}
			return login, nil
		}
		return "", err
	}
	return cookie, nil
}
