package collect

type ParsedRespBody struct {
	Requests []*Request // 先解析出每个页面包含的的url

	Contents []interface{} // 在祖册好的回调被调用的时候会筛选并填充感兴趣的内容
}

// 每一个访问请求都携带着 Url, Cookies 和  对应的响应体处理规则 hook
type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte, *Request) ParsedRespBody
}
