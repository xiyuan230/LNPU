package log

import (
	"MyLNPU/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

func Init() {
	gin.DisableConsoleColor()
	path := conf.GetConfig().Log.Path
	logs, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		Errorf("log文件创建失败... ERROR: %s", err)
		os.Exit(-1)
	}
	gin.DefaultWriter = io.MultiWriter(logs, os.Stdout)
	Println("log初始化成功...")
}
func Println(format string, values ...any) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	f := fmt.Sprintf("[DEV] %s %s\n", now, format)
	fmt.Fprintf(gin.DefaultWriter, f, values...)
}

func Errorf(format string, values ...any) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	f := fmt.Sprintf("[ERR] %s %s\n", now, format)
	fmt.Fprintf(gin.DefaultWriter, f, values...)
}
