package handles

import (
	"MyLNPU/internal/model"
	"MyLNPU/internal/service"
	"github.com/gin-gonic/gin"
)

func GetSystemNotice(c *gin.Context) (any, error) {
	notice, err := service.GetSystemNotice()
	if err != nil {
		return nil, err
	}
	return map[string]any{"notice": notice}, err
}

func UpdateSystemNotice(c *gin.Context) (any, error) {
	var notice model.Notice
	err := c.ShouldBind(&notice)
	if err != nil {
		return nil, err
	}
	err = service.UpdateSystemNotice(&notice)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
