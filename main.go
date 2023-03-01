package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/LinPr/crawler/collect"
	"github.com/LinPr/crawler/log"
	"github.com/LinPr/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var url string = "https://book.douban.com/subject/1007305/"

var url string = "https://google.com"
var rgxPattern = `<div class=[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`

var proxyUrls = []string{"http://192.168.31.67:1080"}

func main() {
	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	// p, err := proxy.NewRoundRobinBalancer(proxyUrls...)
	p, err := proxy.NewConsistentHashBalancer(10, proxyUrls...)
	if err != nil {
		fmt.Printf("proxy.NewRoundRobinBalancer() error: %v\n", err)
		logger.Error("proxy.NewRoundRobinBalancer() failed")
		return
	}
	// var f collect.Fetcher = collect.BaseFetch{}
	var f collect.Fetcher = collect.BrowerFetch{
		Timeout: time.Second * 10,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("Fetch err: %v\n", err)
		logger.Error("read content failed", zap.Error(err))
		return
	}
	logger.Info("get content", zap.Int("len", len(body)))
	fmt.Println(string(body))

	rgx := regexp.MustCompile(rgxPattern)
	matches := rgx.FindAllStringSubmatch(string(body), -1)
	for _, v := range matches {
		fmt.Println(v)
		fmt.Println()
	}

}
