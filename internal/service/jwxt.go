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
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
	"image"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func JwxtLoginWithSSO(openid string) (string, error) {
	cache.Del("lnpu:jwxt:cookie:" + openid)
	user, err := db.GetUserByID(openid)
	if err != nil {
		logger.Errorf("获取用户信息失败... %s", err)
		return "", err
	}
	if user.StudentID == "" || user.SSOPassword == "" {
		return "", errs.ErrUserEmpty
	}
	client, err := loginWithSSO(user.StudentID, user.SSOPassword)
	if err != nil {
		logger.Errorf("统一认证登录失败... %s", err)
		return "", err
	}
	resp, err := client.Get(JwxtLoginUrlWithSSO)
	if err != nil {
		logger.Errorf("教务系统登录失败.... %s", resp)
		return "", err
	}
	defer resp.Body.Close()
	cookie := resp.Request.Header.Get("Cookie")
	cache.Set("lnpu:jwxt:cookie:"+openid, cookie, time.Hour*1)
	return cookie, nil
}

func JwxtLoginWithJwxt(openid string) (string, error) {
	cache.Del("lnpu:jwxt:cookie:" + openid)
	user, err := db.GetUserByID(openid)
	if err != nil {
		logger.Errorf("获取用户信息失败... %s", err)
		return "", err
	}
	if user.StudentID == "" || user.JwxtPassword == "" {
		return "", errs.ErrUserEmpty
	}
	for i := 0; i < 3; i++ {
		code, client, err := utils.GetVerifyCode(JwxtVerifyCodeUrl)
		if err != nil {
			if errors.Is(err, image.ErrFormat) {
				continue
			}
			return "", err
		}
		encode := utils.EncodeByBase64(user.StudentID, user.JwxtPassword)
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
			return "", errs.ErrPasswordWrong
		} else if strings.Contains(string(body), "验证码错误") {
			logger.Println("验证码错误...")
			continue
		}
		cookies := response.Request.Header.Values("Cookie")
		cookie := strings.Join(cookies, "")
		cache.Set("lnpu:jwxt:cookie:"+openid, cookie, time.Hour*1)
		return cookie, nil
	}
	return "", errs.ErrJwxtLoginFailed
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
		logger.Errorf("创建HttpClient对象失败... %s", err)
		return nil, err
	}
	req, _ := http.NewRequest("GET", JwxtStudentInfoUrl, nil)
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("获取学生信息失败... %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Errorf("学生信息页解析失败... %s", err)
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
	logger.Println("获取学生信息[ %s ]成功", stu.Name)
	fmt.Println(cookie)
	return &stu, nil
}

func GetStartDate(openid string) (int64, error) {
	start, err := cache.Get("lnpu:jwxt:startDate")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			doc, err := ParsePage(openid, JwxtCalendarUrl)
			if err != nil {
				return 0, err
			}
			attr, _ := doc.Find("#kbtable tr").Eq(1).Children().Eq(1).Attr("title")
			startDate, err := time.Parse("2006年01月02", attr)
			if err != nil {
				logger.Errorf("学期起始日期格式化失败... %s", err)
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
			doc, err := ParsePage(openid, JwxtScoreUrl)
			if err != nil {
				return nil, err
			}
			listNode := doc.Find(".Nsb_r_list").Children().Children()
			length := listNode.Length()
			if length == 0 {
				cache.Del("lnpu:jwxt:cookie:" + openid)
				return nil, errs.ErrCookieExpire
			}
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

func GetCourseTable(openid string) (*[]model.Course, error) {
	data, err := cache.Get("lnpu:jwxt:course:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			doc, err := ParsePage(openid, JwxtCourseUrl)
			if err != nil {
				return nil, err
			}
			var courseList []model.Course
			courseNode := doc.Find("#kbtable tr")
			if courseNode.Length() == 0 {
				cache.Del("lnpu:jwxt:cookie:" + openid)
				return nil, errs.ErrCookieExpire
			}
			courseNode.Find("th").Remove()
			courseNode.Find("input").Remove()
			courseNode.Find(".kbcontent1").Remove()
			courseNode.Find(".sykb1").Remove()
			courseNode.Find(".sykb2").Remove()
			for i := 1; i <= 6; i++ {
				section := i
				courseNode.Eq(i).Children().Each(func(j int, selection *goquery.Selection) {
					week := j + 1
					if ok, _ := regexp.MatchString("[\u4e00-\u9fa5]", selection.Text()); ok {
						kbcontent := selection.Find(".kbcontent").Eq(0)
						count := kbcontent.Children().Length()
						if count > 7 {
							for k := 0; k < count/7; k++ {
								course := model.Course{}
								course.CourseName = kbcontent.Contents().Eq(0 + k*10).Text()
								course.Address = kbcontent.Find("[title='教室']").Eq(k).Text()
								course.Teacher = kbcontent.Find("[title='老师']").Eq(k).Text()
								course.WeekList = utils.CourseWeekListHandle(kbcontent.Find("[title='周次(节次)']").Eq(k).Text())
								course.Week = week
								course.Sections = section
								courseList = append(courseList, course)
							}
						} else {
							course := model.Course{}
							course.CourseName = kbcontent.Contents().Eq(0).Text()
							course.Address = kbcontent.Find("[title='教室']").Eq(0).Text()
							course.Teacher = kbcontent.Find("[title='老师']").Eq(0).Text()
							course.WeekList = utils.CourseWeekListHandle(kbcontent.Find("[title='周次(节次)']").Eq(0).Text())
							course.Week = week
							course.Sections = section
							courseList = append(courseList, course)
						}
					}
				})
			}
			marshal, _ := json.Marshal(courseList)
			cache.Set("lnpu:jwxt:course:"+openid, marshal, time.Hour*12)
			return &courseList, nil
		}
		return nil, err
	}
	var courseTable []model.Course
	json.Unmarshal([]byte(data), &courseTable)
	return &courseTable, nil
}

func GetTrainingTable(openid string) (*[]model.CourseDetail, error) {
	data, err := cache.Get("lnpu:jwxt:training:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			doc, err := ParsePage(openid, JwxtTrainingTableUrl)
			if err != nil {
				return nil, err
			}
			doc.Find("[colspan='7']").Parent().Remove()
			doc.Find("TR").Last().Remove()
			course := doc.Find("#mxh tr")
			index := 0
			var table []model.CourseDetail
			for i := 2; i < course.Length(); i++ {
				detail := model.CourseDetail{}
				if course.Eq(i).Children().Length() == 14 {
					index = i
					detail.Type = utils.TrimSpace(course.Eq(i).Children().Eq(0).Text())
					detail.CourseName = utils.TrimSpace(course.Eq(i).Children().Eq(3).Text())
					if utils.TrimSpace(course.Eq(i).Children().Eq(4).Text()) == "" {
						detail.Status = "未开始"
					} else {
						detail.Status = utils.TrimSpace(course.Eq(i).Children().Eq(4).Text())
					}
					detail.Attribute = utils.TrimSpace(course.Eq(i).Children().Eq(6).Text())
					detail.Credit = utils.TrimSpace(course.Eq(i).Children().Eq(7).Text())
					detail.CreditHours = utils.TrimSpace(course.Eq(i).Children().Eq(12).Text())
					var term string
					switch utils.TrimSpace(course.Eq(i).Children().Eq(12).Text()) {
					case "1":
						term = "大一上"
					case "2":
						term = "大一下"
					case "3":
						term = "大二上"
					case "4":
						term = "大二下"
					case "5":
						term = "大三上"
					case "6":
						term = "大三下"
					case "7":
						term = "大四上"
					case "8":
						term = "大四下"
					}
					detail.Term = term
				} else {
					detail.Type = utils.TrimSpace(course.Eq(index).Children().Eq(0).Text())
					detail.CourseName = utils.TrimSpace(course.Eq(i).Children().Eq(2).Text())
					if utils.TrimSpace(course.Eq(i).Children().Eq(3).Text()) == "" {
						detail.Status = "未开始"
					} else {
						detail.Status = utils.TrimSpace(course.Eq(i).Children().Eq(3).Text())
					}
					detail.Attribute = utils.TrimSpace(course.Eq(i).Children().Eq(5).Text())
					detail.Credit = utils.TrimSpace(course.Eq(i).Children().Eq(6).Text())
					detail.CreditHours = utils.TrimSpace(course.Eq(i).Children().Eq(11).Text())
					var term string
					switch utils.TrimSpace(course.Eq(i).Children().Eq(12).Text()) {
					case "1":
						term = "大一上"
					case "2":
						term = "大一下"
					case "3":
						term = "大二上"
					case "4":
						term = "大二下"
					case "5":
						term = "大三上"
					case "6":
						term = "大三下"
					case "7":
						term = "大四上"
					case "8":
						term = "大四下"
					}
					detail.Term = term
				}
				table = append(table, detail)
			}
			marshal, _ := json.Marshal(table)
			cache.Set("lnpu:jwxt:training:"+openid, marshal, time.Hour*24)
			return &table, nil
		}
		return nil, err
	}
	var table []model.CourseDetail
	json.Unmarshal([]byte(data), &table)
	return &table, nil
}

// UpdateCookie 更新cookie
func UpdateCookie(openid string) (string, error) {
	cookie, err := cache.Get("lnpu:jwxt:cookie:" + openid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			login, err := JwxtLoginWithSSO(openid)
			if err != nil {
				if errors.Is(err, errs.ErrUserEmpty) {
					login, err := JwxtLoginWithJwxt(openid)
					if err != nil {
						return "", err
					}
					return login, nil
				}
				return "", err
			}
			return login, nil
		}
		return "", err
	}
	return cookie, nil
}

func ParsePage(openid, url string) (*goquery.Document, error) {
	cookie, err := UpdateCookie(openid)
	if err != nil {
		return nil, err
	}
	client, err := utils.NewHttpClient()
	if err != nil {
		logger.Errorf("创建HttpClient失败... %s", err)
		return nil, err
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
