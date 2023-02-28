package main

import (
	"fmt"
	"github.com/LinPr/crawler/collect"
	"github.com/LinPr/crawler/proxy"
	"regexp"
	"time"
)

// var url string = "https://book.douban.com/subject/1007305/"

var url string = "https://google.com"
var rgxPattern = `<div class=[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`

var proxyUrls = []string{"http://192.168.31.67:1080"}

func main() {

	p, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	if err != nil {
		fmt.Printf("proxy.RoundRobinProxySwitcher*() error: %v\n", err)
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
		return
	}

	fmt.Println(string(body))

	rgx := regexp.MustCompile(rgxPattern)
	matches := rgx.FindAllStringSubmatch(string(body), -1)
	for _, v := range matches {
		fmt.Println(v)
		fmt.Println()
	}

}
