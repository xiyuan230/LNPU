package api

import (
	"MyLNPU/api/common"
	"MyLNPU/api/handles"
	"MyLNPU/api/middlewares"
	"MyLNPU/internal/errs"
	"MyLNPU/internal/log"
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	auth(e.Group("/auth"))   //用户路由
	admin(e.Group("/admin")) //管理路由
	other(e.Group("/other"))
	jwxt(e.Group("/jwxt"))
	experiment(e.Group("/experiment"))
	log.Println("router初始化成功...")
}

func auth(r *gin.RouterGroup) {
	r.Use(middlewares.AuthorizationWithToken)
	r.GET("/login", Wrapper(handles.Login))
	r.GET("/status", Wrapper(handles.CheckTokenStatus))
}

func admin(r *gin.RouterGroup) {
	r.Use(middlewares.AuthorizationWithAdmin)
	r.POST("/notice", Wrapper(handles.UpdateSystemNotice))
}

func other(r *gin.RouterGroup) {
	r.GET("/notice", Wrapper(handles.GetSystemNotice)) //获取公告
}

func jwxt(r *gin.RouterGroup) {
	r.Use(middlewares.AuthorizationWithToken)
	r.GET("/login", Wrapper(handles.JwxtLogin))
	r.GET("/startDate", Wrapper(handles.GetStartDate))
	r.GET("/score", Wrapper(handles.GetJwxtScore))
	r.GET("/course", Wrapper(handles.GetCourseTable))
	r.GET("/training", Wrapper(handles.GetTrainingTable))
}
func experiment(r *gin.RouterGroup) {
	r.GET("/table", Wrapper(handles.GetExpTable))
}

// WrapperHandle 全局统一错误处理
type WrapperHandle func(c *gin.Context) (interface{}, error)

func Wrapper(handle WrapperHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handle(c)
		if err != nil {
			switch err {
			case errs.ErrPasswordWrong:
				common.ErrorStrResp(c, 401, "账号或密码错误")
				return
			case errs.ErrUserEmpty:
				common.ErrorStrResp(c, 401, "还未绑定身份信息")
				return
			case errs.ErrParamMiss:
				common.ErrorStrResp(c, 401, "请求参数错误")
				return
			case errs.ErrCookieExpire:
				common.ErrorStrResp(c, 500, "Cookie失效，请刷新")
				return
			default:
				common.ErrorResp(c, 500, err)
				return
			}
		}
		common.SuccessResp(c, data)
	}
}
