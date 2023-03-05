// schedule 其实是个循环流水线，以下是流程图
//
// 		rootPaths
// 			↓
// 			↓ 		  workerC						outputC
// =====> reqQueue ===========>> ParseFunc(req)  =============>>+
// ↑									 						↓
// ↑															↓
// + <<================================================== HandleResult() - item ==> 存储
// 					requestC

package pipeline

import (
	"time"

	"github.com/LinPr/crawler/collect"
	"go.uber.org/zap"
)

type PipelineEngine struct {
	// 对象
	Seeds       []*collect.Request
	Fetcher     collect.Fetcher
	WorkerCount int
	Logger      *zap.Logger
	WaitTime    time.Duration

	// 中转管道
	RequestC chan *collect.Request
	WorkerC  chan *collect.Request
	OutputC  chan collect.ParsedRespBody
}

func (pe *PipelineEngine) Run() {

	go pe.HandleReq()

	for i := 0; i < pe.WorkerCount; i++ {
		go pe.CreateWorker(i)
	}

	// 主协程
	pe.HandleRespBody()
}

func (pe *PipelineEngine) HandleReq() {
	var crawlReqQueue = pe.Seeds

	for {
		var req *collect.Request
		var ch chan *collect.Request

		if len(crawlReqQueue) > 0 {
			req = crawlReqQueue[0]
			ch = pe.WorkerC
		}

		select {
		// 将新到 req 放入请求队列末尾
		case newReq := <-pe.RequestC:
			crawlReqQueue = append(crawlReqQueue, newReq)

		// 如果前面的if不执行，这里就会阻塞，因为ch是nil
		case ch <- req:
			crawlReqQueue = crawlReqQueue[1:]
			time.Sleep(pe.WaitTime)

		}

	}

}

func (pe *PipelineEngine) CreateWorker(id int) {
	// fmt.Printf("\ncreate worker id %v\n", id)
	for {
		req := <-pe.WorkerC

		body, err := pe.Fetcher.Get(req)
		if err != nil {
			pe.Logger.Error("s.Fetcher.Get(req)", zap.Error(err))
			continue
		}
		parsedRespBody := req.ParseFunc(body, req)

		pe.OutputC <- parsedRespBody

	}
}

func (pe *PipelineEngine) HandleRespBody() {

	// 使用range for 来循环读channel，只要channel不被关闭，就会一直循环
	for parsedRespBody := range pe.OutputC {
		for _, req := range parsedRespBody.Requests {
			pe.RequestC <- req
		}

		for _, content := range parsedRespBody.Contents {
			pe.Logger.Sugar().Info("get content", content.(string))
		}
	}
}
