package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var config ApplicationConfig

type ApplicationConfig struct {
	ApplicationName string
	Server          ServerConfig
	Log             LogConfig
	Redis           RedisConfig
	MySql           MysqlConfig
	Proxy           ProxyConfig
}

type ServerConfig struct {
	Port int
}
type LogConfig struct {
	Path string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MysqlConfig struct {
	Host         string
	Port         string
	UserName     string
	UserPassword string
	Database     string
}

type ProxyConfig struct {
	EnableProxy bool
	ProxyUrl    string
}

func Init() {
	viper.SetConfigName("application")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./conf")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("读取配置文件失败... ERROR: %s\n", err)
		os.Exit(-1)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("配置文件解析失败... ERROR: %s\n", err)
		os.Exit(-1)
	}
	viper.WatchConfig()
	printBanner()
	log.Println("配置文件读取成功...")
}

func GetConfig() ApplicationConfig {
	return config
}

func printBanner() {
	fmt.Println("                 __  __           _       _   _   ____    _   _ \n                |  \\/  |  _   _  | |     | \\ | | |  _ \\  | | | |\n                | |\\/| | | | | | | |     |  \\| | | |_) | | | | |\n                | |  | | | |_| | | |___  | |\\  | |  __/  | |_| |\n                |_|  |_|  \\__, | |_____| |_| \\_| |_|      \\___/ \n                          |___/                                 \n")
}
