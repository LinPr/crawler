package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"regexp"
)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)
	if err != nil {
		fmt.Printf("Fetch err: %v\n", err)
		return
	}

	// fmt.Println(string(body))
	rgx := regexp.MustCompile(`<div class=[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
	matches := rgx.FindAllStringSubmatch(string(body), -1)
	for _, v := range matches {
		fmt.Println(v)
		fmt.Println()
	}

}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("http.Get(url) err: %v\n", err)
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)

	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)

	// fmt.Printf("body: %v\n", string(body))
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	//检测html页面编码使用peek
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("r.Peek(1024) err: %v\n", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
