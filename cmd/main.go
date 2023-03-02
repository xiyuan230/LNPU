package main

import (
	"MyLNPU/api"
	"MyLNPU/conf"
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	conf.Init()  //配置初始化
	log.Init()   //日志初始化
	api.Init(r)  //路由初始化
	cache.Init() //Redis初始化
	db.Init()    //Mysql初始化
	serverPort := ":" + strconv.Itoa(conf.GetConfig().Server.Port)
	log.Println("MyLNPU启动成功... Port: %d", conf.GetConfig().Server.Port)
	if err := r.Run(serverPort); err != nil {
		log.Errorf(err.Error())
	}
}
