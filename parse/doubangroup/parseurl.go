package doubangroup

import (
	"fmt"
	"os"
	"regexp"

	"github.com/LinPr/crawler/collect"
)

const topicUrlRgx string = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseUrl(respBody []byte, req *collect.Request) collect.ParsedRespBody {
	f, err := os.OpenFile("./demo.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 644)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	fmt.Fprintf(f, "%v", string(respBody))

	rgx := regexp.MustCompile(topicUrlRgx)

	matches := rgx.FindAllStringSubmatch(string(respBody), -1)

	var result = collect.ParsedRespBody{}

	for _, m := range matches {
		u := string(m[1])
		fmt.Printf("u: %v\n", u)

		// 为每个content条目注册的回调函数，为了传参，需要外面再包装一层（参考C++的placeholder机制）
		// f := func(respBody []byte, req *collect.Request) collect.ParsedRespBody {
		// 	prb := getNeedContent(respBody, req)
		// 	return prb
		// }

		req := collect.Request{
			Url:       u,
			Cookie:    req.Cookie,
			ParseFunc: getNeedContent,
		}
		result.Requests = append(result.Requests, &req)
	}
	return result
}

const contentRe string = `[\s\S]*?一室一厅[\s\S]*?`

func getNeedContent(respBody []byte, req *collect.Request) collect.ParsedRespBody {
	cr := regexp.MustCompile(contentRe)
	ok := cr.Match(respBody)
	fmt.Printf("ok: %v\n", ok)
	if !ok {
		return collect.ParsedRespBody{
			Contents: []interface{}{},
		}
	}

	// 如果找到感兴趣的内容， 就将指向该内容的url返回
	prb := collect.ParsedRespBody{
		Contents: []interface{}{req.Url},
	}
	return prb
}
