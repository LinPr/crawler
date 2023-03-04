# crawler
A web crawler project

### 实现的功能
1. 最基本的访问网站，发送 http request 请求。获取 http response 的能力
2. 对有些非 UTF-8 编码的 response 进行统一的转码处理
3. 使用正则表达式 regexp 对返回的 response 中的有效数据进行数据提取过滤
4. 模拟浏览器访问网站来绕过网站反爬虫机制，包括设置 http request 的头部 User-Agent 字段, proxy, cookies， ip代理池
5. http.Client对象默认会从环境变量读取 HTTP_PROXY, HTTPS_PROXY, NO_PROXY, 其作用之一为反爬，另一个作用为访问外网
6. 可以通过设置clinet的transport成员来修改代理规则，并且使用分别使用轮询和一致性哈希算法作为负载均衡算法
7. 对不同类型网站需要设置不同的遍历算法，例如，豆瓣这种翻页式网站需要分析URL手动构建URL队列，对于树形超链接则需要使用深度优先或者广度优先获取所有URL，同时防止链接环路需要先利用广度优先拓扑排序算法将URL遍历一遍后放到URL队列中，然后再按顺序进行访问，这里一般在实际项目中采用广度优先，因为深度优先采用递归实现会造成递归栈过深的问题
