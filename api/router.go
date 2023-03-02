package api

import (
	"MyLNPU/api/handles"
	"MyLNPU/api/middlewares"
	"MyLNPU/internal/log"
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	auth(e.Group("/auth"))   //用户路由
	admin(e.Group("/admin")) //管理路由
	common(e.Group("/common"))

	e.GET("/test", handles.TestSSO)
	log.Println("router初始化成功...")
}

func auth(r *gin.RouterGroup) {
	r.Use(middlewares.AuthorizationWithToken)
	r.GET("/login", handles.Login)
	r.GET("/status", handles.CheckTokenStatus)
}

func admin(r *gin.RouterGroup) {
	r.Use(middlewares.AuthorizationWithAdmin)
	r.POST("/notice", handles.UpdateSystemNotice)
}

func common(r *gin.RouterGroup) {
	r.GET("/notice", handles.GetSystemNotice) //获取公告
}
