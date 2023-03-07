package main

import "MyLNPU/cmd"

func main() {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.Default()
	//conf.Init()   //配置初始化
	//logger.Init() //日志初始化
	//api.Init(r)   //路由初始化
	//cache.Init()  //Redis初始化
	//db.Init()     //Mysql初始化
	//serverPort := ":" + strconv.Itoa(conf.GetConfig().Server.Port)
	//logger.Println("MyLNPU启动成功... Port: %d", conf.GetConfig().Server.Port)
	//if err := r.Run(serverPort); err != nil {
	//	logger.Errorf(err.Error())
	//}
	cmd.Execute()
}
