package service

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/errs"
	"MyLNPU/internal/log"
	"MyLNPU/internal/model"
	"MyLNPU/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func JwxtLogin(openid string) (string, error) {
	cache.Del("lnpu:jwxt:cookie:" + openid)
	user, err := db.GetUserByID(openid)
	if err != nil {
		log.Errorf("获取用户信息失败... %s", err)
		return "", err
	}
	if user.StudentID == "" || user.SSOPassword == "" {
		return "", errs.ErrUserEmpty
	}
	client, err := SSOLogin(user.StudentID, user.SSOPassword)
	if err != nil {
		log.Errorf("统一认证登录失败... %s", err)
		return "", err
	}
	resp, err := client.Get(JwxtLoginUrlWithSSO)
	if err != nil {
		log.Errorf("教务系统登录失败.... %s", resp)
		return "", err
	}
	defer resp.Body.Close()
	cookie := resp.Request.Header.Get("Cookie")
	cache.Set("lnpu:jwxt:cookie:"+openid, cookie, time.Hour*1)
	return cookie, nil
}

// GetStudentInfo 获取学生信息
func GetStudentInfo(openid string) (*model.Student, error) {
	stu := model.Student{}
	cookie, err := UpdateCookie(openid)
	if err != nil {
		return nil, err
	}
	client, err := utils.NewHttpClient()
	if err != nil {
		log.Errorf("创建HttpClient对象失败... %s", err)
		return nil, err
	}
	req, _ := http.NewRequest("GET", JwxtStudentInfoUrl, nil)
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("获取学生信息失败... %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("学生信息页解析失败... %s", err)
		return nil, err
	}
	infoNode := doc.Find(".middletopdwxxcont")
	if infoNode.Length() == 0 {
		return nil, errs.ErrCookieExpire
	}
	stu.Name = infoNode.Eq(1).Text()
	stu.StudentID = infoNode.Eq(2).Text()
	stu.College = infoNode.Eq(3).Text()
	stu.Major = infoNode.Eq(4).Text()
	stu.Class = infoNode.Eq(5).Text()
	log.Println("获取学生信息[ %s ]成功", stu.Name)
	fmt.Println(cookie)
	return &stu, nil
}

func GetStartDate(openid string) (int64, error) {
	start, err := cache.Get("lnpu:jwxt:startDate")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			cookie, err := UpdateCookie(openid)
			if err != nil {
				return 0, err
			}
			client, err := utils.NewHttpClient()
			if err != nil {
				log.Errorf("创建HttpClient失败... %s", err)
				return 0, err
			}
			req, _ := http.NewRequest("GET", JwxtCalendarUrl, nil)
			req.Header.Add("Cookie", cookie)
			resp, err := client.Do(req)
			if err != nil {
				return 0, err
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return 0, err
			}
			attr, _ := doc.Find("#kbtable tr").Eq(1).Children().Eq(1).Attr("title")
			startDate, err := time.Parse("2006年01月02", attr)
			if err != nil {
				log.Errorf("学期起始日期格式化失败... %s", err)
				return 0, err
			}
			cache.Set("lnpu:jwxt:startDate", startDate.Unix(), time.Hour*24)
			return startDate.Unix(), nil
		}
		return 0, err
	}
	parseInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}

func GetJwxtScore(openid string) (*model.ScoreResult, error) {
	data, err := cache.Get("lnpu:jwxt:score:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			cookie, err := UpdateCookie(openid)
			if err != nil {
				return nil, err
			}
			client, err := utils.NewHttpClient()
			if err != nil {
				return nil, err
			}
			req, _ := http.NewRequest("GET", JwxtScoreUrl, nil)
			req.Header.Add("Cookie", cookie)
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return nil, err
			}
			listNode := doc.Find(".Nsb_r_list").Children().Children()
			length := listNode.Length()
			scoreList := make([]model.Score, length-1)
			for i := 1; i < length; i++ {
				score := model.Score{}
				score.Term = strings.TrimSpace(listNode.Eq(i).Children().Eq(1).Text())
				score.ClassName = strings.TrimSpace(listNode.Eq(i).Children().Eq(3).Text())
				score.Score = strings.TrimSpace(listNode.Eq(i).Children().Eq(5).Text())
				score.Credits = strings.TrimSpace(listNode.Eq(i).Children().Eq(7).Text())
				score.GPA = strings.TrimSpace(listNode.Eq(i).Children().Eq(9).Text())
				score.Pattern = strings.TrimSpace(listNode.Eq(i).Children().Eq(14).Text())
				scoreList[i-1] = score
			}
			result := model.ScoreResult{}
			result.ScoreList = scoreList
			info := doc.Find(".Nsb_r_list").Parent()
			info.Children().Remove()
			text := info.Text()
			text = strings.ReplaceAll(text, "\n", "")
			text = strings.ReplaceAll(text, "：", ":")
			compile := regexp.MustCompile("[\u4e00-\u9fa5]+:[0-9]+.*[0-9]*")
			subMatch := compile.FindStringSubmatch(text)
			splitTmp := strings.Split(subMatch[0], " ")
			var split []string
			for i := 0; i < len(splitTmp); i++ {
				if ok, _ := regexp.MatchString("[\u4e00-\u9fa5]+:[0-9]+.*[0-9]*", splitTmp[i]); ok {
					split = append(split, splitTmp[i])
				}
			}
			result.CourseCount = utils.ScoreStrHandle(split[0])
			result.TotalCredit = utils.ScoreStrHandle(split[1])
			result.AverageCreditPoint = utils.ScoreStrHandle(split[2])
			result.AverageGrade = utils.ScoreStrHandle(split[3])
			result.Rank = utils.ScoreStrHandle(split[4])
			marshal, _ := json.Marshal(result)
			cache.Set("lnpu:jwxt:score:"+openid, marshal, time.Hour*2)
			return &result, nil
		}
		return nil, err
	}
	var scoreResult model.ScoreResult
	json.Unmarshal([]byte(data), &scoreResult)
	return &scoreResult, nil
}

// UpdateCookie 更新cookie
func UpdateCookie(openid string) (string, error) {
	cookie, err := cache.Get("lnpu:jwxt:cookie:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			login, err := JwxtLogin(openid)
			if err != nil {
				return "", err
			}
			return login, nil
		}
		return "", err
	}
	return cookie, nil
}
