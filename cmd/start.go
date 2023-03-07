package cmd

import (
	"MyLNPU/api"
	"MyLNPU/conf"
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/logger"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

var daemon bool

func Execute() {
	flag.BoolVar(&daemon, "d", false, "Enable daemon process")
	flag.Parse()
	conf.Init() //配置初始化
	if daemon {
		cmd, err := startProc()
		if err != nil {
			log.Println("启动子进程失败...", err)
			return
		}
		if cmd != nil {
			_, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetConfig().Server.Port))
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("启动成功... PID: ", cmd.Process.Pid)
			return
		}
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	logger.Init() //日志初始化
	api.Init(r)   //路由初始化
	cache.Init()  //Redis初始化
	db.Init()     //Mysql初始化
	serverPort := ":" + strconv.Itoa(conf.GetConfig().Server.Port)
	logger.Println("MyLNPU启动成功... Port: %d", conf.GetConfig().Server.Port)
	if err := r.Run(serverPort); err != nil {
		log.Println(err)
	}
}

func startProc() (*exec.Cmd, error) {
	envName := "FLAG_DAEMON"
	envValue := "SUB_PROC"
	val := os.Getenv(envName)
	if val == envValue {
		return nil, nil
	}
	cmd := &exec.Cmd{Path: os.Args[0], Args: os.Args, Env: os.Environ()}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))
	file, err := os.OpenFile("/var/log/lnpu.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(os.Getpid(), ": 打开日志文件失败...", err)
		return nil, err
	}
	cmd.Stdout = file
	cmd.Stderr = file

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
