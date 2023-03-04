package main

import (
	"fmt"
	// "regexp"
	"time"

	"github.com/LinPr/crawler/collect"
	"github.com/LinPr/crawler/log"
	"github.com/LinPr/crawler/parse/doubangroup"
	"github.com/LinPr/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var url string = "https://book.douban.com/subject/1007305/"

var url = "https://google.com"
var rgxPattern = `<div class=[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`
var proxyUrls = []string{"http://192.168.31.67:1080"}
var doubanCookie = `dbcl2="268136649:cAIO0noyRks"; bid=Ps7UbSYL1Bo; ck=7_wQ; __gads=ID=0c2a830c2f1225d0-22e828845eda00bd:T=1677859864:RT=1677859864:S=ALNI_MagcORlvrA3lKBJlTF6eg--LT9FmA; push_noty_num=0; push_doumail_num=0; __utmc=30149280; __utmv=30149280.26813; _pk_ref.100001.8cb4=["","",1677905494,"https://cn.bing.com/"]; _pk_ses.100001.8cb4=*; ap_v=0,6.0; __yadk_uid=mlWttd31znPUb98VivTD7pkUIWdgVnJ3; __utma=30149280.1558330668.1677859864.1677859864.1677905496.2; __utmz=30149280.1677905496.2.2.utmcsr=cn.bing.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __gpi=UID=00000bd126631ba6:T=1677859864:RT=1677905502:S=ALNI_MbtdxhXDFfTJHvzn3DKgYMZ4597qA; frodotk_db="b45b01db7fa8fe581d1105927efafb08"; douban-fav-remind=1; ct=y; _pk_id.100001.8cb4=e5ecf39d85aaafb9.1677859860.2.1677905899.1677859875.; __utmb=30149280.237.5.1677905900199`

func main() {
	// log init
	// plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	// defer c.Close()
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)

	// proxy init
	proxyUrls = nil
	// p, err := proxy.NewRoundRobinBalancer(proxyUrls...)
	p, err := proxy.NewConsistentHashBalancer(10, proxyUrls...)
	if err != nil {
		fmt.Printf("proxy.NewRoundRobinBalancer() error: %v\n", err)
		logger.Error("proxy.NewRoundRobinBalancer() failed")
		// 不需要 return， 会使用默认
	}

	// init crawl request queue
	var crawlReqQueue []*collect.Request
	for i := 0; i < 1; i++ {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		req := collect.Request{
			Url:       str,
			Cookie:    doubanCookie,
			ParseFunc: doubangroup.ParseUrl,
		}
		crawlReqQueue = append(crawlReqQueue, &req)
	}

	// fetcher
	// var f collect.Fetcher = collect.BaseFetch{}
	var f collect.Fetcher = collect.BrowerFetch{
		Timeout: time.Second * 10,
		Proxy:   p,
	}

	seen := make(map[*collect.Request]bool) //防止遍历节点存在环
	// 用一个for循环解决层序遍历，每循环一轮，遍历个页面
	for len(crawlReqQueue) > 0 {
		tmpCrawlReqQueue := crawlReqQueue
		crawlReqQueue = nil
		for _, req := range tmpCrawlReqQueue {
			body, err := f.Get(req)
			time.Sleep(time.Second * 1)
			if err != nil {
				logger.Error("f.Get(req) failed", zap.Error(err))
				continue
			}

			parsedReqBody := req.ParseFunc(body, req)
			for _, content := range parsedReqBody.Contents {
				logger.Info("parsedReqbody", zap.String("get url:", content.(string)))
			}
			if !seen[req] {
				seen[req] = true
				crawlReqQueue = append(crawlReqQueue, parsedReqBody.Requests...)
			}

		}
	}

	// —————————————这些不用，就注释掉了—————————————————
	// body, err := f.Get(url)
	// if err != nil {
	// 	fmt.Printf("Fetch err: %v\n", err)
	// 	logger.Error("read content failed", zap.Error(err))
	// 	return
	// }
	// logger.Info("get content", zap.Int("len", len(body)))
	// fmt.Println(string(body))

	// rgx := regexp.MustCompile(rgxPattern)
	// matches := rgx.FindAllStringSubmatch(string(body), -1)
	// for _, v := range matches {
	// 	fmt.Println(v)
	// 	fmt.Println()
	// }

}
