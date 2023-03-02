package utils

import (
	"MyLNPU/conf"
	"MyLNPU/internal/log"
	"golang.org/x/net/proxy"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func NewHttpClient() (*http.Client, error) {
	var client http.Client
	config := conf.GetConfig().Proxy
	Enable := config.EnableProxy
	ProxyUrl := config.ProxyUrl
	if Enable {
		socks5, err := proxy.SOCKS5("tcp", ProxyUrl, nil, proxy.Direct)
		if err != nil {
			log.Errorf("无法连接到代理服务器... %s", err)
			return nil, err
		}
		jar, _ := cookiejar.New(nil)
		httpTransport := http.Transport{DialContext: socks5.(proxy.ContextDialer).DialContext}
		client = http.Client{
			Transport: &httpTransport,
			Jar:       jar,
			Timeout:   time.Second * 15,
		}
	} else {
		jar, _ := cookiejar.New(nil)
		client = http.Client{
			Transport: nil,
			Jar:       jar,
			Timeout:   time.Second * 15,
		}
	}
	return &client, nil
}
