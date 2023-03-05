package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LinPr/crawler/proxy"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(r *Request) ([]byte, error)
}

type BaseFetch struct {
}

// func (BaseFetch) Get(r *Request) ([]byte, error) {
func (BaseFetch) Get(url string) ([]byte, error) {
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

type BrowerFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
	Logger  *zap.Logger
}

// func (b BrowerFetch) Get(url string) ([]byte, error) {
func (b BrowerFetch) Get(r *Request) ([]byte, error) {
	client := http.Client{
		Timeout: b.Timeout,
	}
	if b.Proxy != nil {
		transport, _ := http.DefaultTransport.(*http.Transport) // 类型断言，为了使用concrete type数据的成员和方法
		transport.Proxy = b.Proxy
		client.Transport = transport // 配置client对象的transport为自定义
	}
	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		// fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, fmt.Errorf("http.NewRequest(): %v", err)
	}

	// 设置HTTP请求头部字段
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57")
	if len(r.Cookie) > 0 {
		req.Header.Set("Cookie", r.Cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		b.Logger.Error("client.Do(req)", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()
	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)

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
