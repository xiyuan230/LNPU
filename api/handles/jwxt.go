package handles

import (
	"MyLNPU/internal/errs"
	"MyLNPU/internal/log"
	"MyLNPU/internal/service"
	"MyLNPU/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func JwxtLogin(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	_, err = service.JwxtLogin(openid)
	if err != nil {
		if errors.Is(err, errs.ErrPasswordWrong) {
			return nil, err
		}
		return nil, err
	}
	stu, err := service.GetStudentInfo(openid)
	if err != nil {
		return nil, err
	}
	return map[string]any{"student_info": stu}, err
}

func GetStartDate(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	startDate, err := service.GetStartDate(openid)
	if err != nil {
		log.Errorf("获取学期起始日期失败... %s", err)
		return nil, err
	}
	log.Println("获取学期起始日期成功 [%s]", openid)
	return map[string]any{"start_date": startDate}, nil
}

func GetJwxtScore(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	score, err := service.GetJwxtScore(openid)
	if err != nil {
		log.Errorf("获取成绩信息失败... %s", err)
		return nil, err
	}
	log.Println("获取成绩信息成功 [%s]", openid)
	return map[string]any{"score": score}, nil
}

func GetCourseTable(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	courseTable, err := service.GetCourseTable(openid)
	if err != nil {
		log.Errorf("获取课表失败... %s", err)
		return nil, err
	}
	log.Println("获取课表信息成功 [%s]", openid)
	return map[string]any{"course_table": courseTable}, nil
}

func GetTrainingTable(c *gin.Context) (any, error) {
	token := c.GetHeader("Authorization")
	openid, err := utils.JWTParseToken(token)
	if err != nil {
		return nil, err
	}
	table, err := service.GetTrainingTable(openid)
	if err != nil {
		log.Errorf("获取培养方案信息失败... %s", err)
		return nil, err
	}
	log.Println("获取培养方案信息成功 [%s]", openid)
	return map[string]any{"training_table": table}, nil
}
