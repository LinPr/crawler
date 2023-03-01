# crawler
A web crawler project

### 实现的功能
1. 最基本的访问网站，发送 http request 请求。获取 http response 的能力
2. 对有些非 UTF-8 编码的 response 进行统一的转码处理
3. 使用正则表达式 regexp 对返回的 response 中的有效数据进行数据提取过滤
4. 模拟浏览器访问网站来绕过网站反爬虫机制，包括设置 http request 的头部 User-Agent 字段, proxy, cookies
5. http.Client对象默认会从环境变量读取 HTTP_PROXY, HTTPS_PROXY, NO_PROXY
6. 可以通过设置clinet的transport成员来修改代理规则，并且使用分别使用轮询和一致性哈希算法作为负载均衡算法
