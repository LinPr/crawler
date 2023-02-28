package main

import (
	"fmt"
	"github.com/LinPr/crawler/collect"
	"regexp"
	"time"
)

var url string = "https://book.douban.com/subject/1007305/"
var rgxPattern = `<div class=[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`

func main() {

	// var f collect.Fetcher = collect.BaseFetch{}
	var f collect.Fetcher = collect.BrowerFetch{
		Timeout: time.Second * 3,
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
