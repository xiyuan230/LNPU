package handles

import (
	"MyLNPU/api/common"
	"MyLNPU/internal/model"
	"MyLNPU/internal/service"
	"github.com/gin-gonic/gin"
)

func GetSystemNotice(c *gin.Context) {
	notice, err := service.GetSystemNotice()
	if err != nil {
		common.ErrorResp(c, 500, err)
	}
	common.SuccessResp(c, map[string]any{"notice": notice})
}

func UpdateSystemNotice(c *gin.Context) {
	var notice model.Notice
	err := c.ShouldBind(&notice)
	if err != nil {
		common.ErrorResp(c, 401, err)
		return
	}
	err = service.UpdateSystemNotice(&notice)
	if err != nil {
		common.ErrorResp(c, 500, err)
		return
	}
	common.SuccessResp(c, nil)
}
