package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LinPr/crawler/proxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BaseFetch struct {
}

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
}

func (b BrowerFetch) Get(url string) ([]byte, error) {
	client := http.Client{
		Timeout: b.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, fmt.Errorf("http.NewRequest(): %v\n", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
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
