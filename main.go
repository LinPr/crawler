package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %v", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(body))

	count := strings.Count(string(body), "<a") // 几乎等价 bytes.Count()
	fmt.Printf("links count: %v\n", count)

	exist := strings.Contains(string(body), "中国")
	fmt.Printf("contains exist: %v\n", exist) // 几乎等价  bytes.Contains()

}
