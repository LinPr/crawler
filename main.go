package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	url := "https://www.jd.com"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("http.Get(url) err: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("io.ReadAll(resp.Body) err: %v\n", err)
		return
	}
	// fmt.Printf("body: %v\n", string(body))

	numLinks := strings.Count(string(body), "<a")
	fmt.Printf("homepage has %d sublinks!\n", numLinks)

	numLinks = bytes.Count(body, []byte("<a"))
	fmt.Printf("homepage has %d sublinks!\n", numLinks)

	exist := strings.Contains(string(body), "<a")
	fmt.Printf("homepage exist link? %v\n", exist)

	exist = bytes.Contains(body, []byte("<a"))
	fmt.Printf("homepage exist link? %v\n", exist)

}
