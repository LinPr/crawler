package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyURls []*url.URL
	index     uint32
}

// 这是一个闭包，
// 在go语言中函数作为一等公民，函数可以当做任何参数或者返回值传来传去（类似C++中的std::function，或者仿函数对象）
// 使用方法为，根据这个闭包生成一个对象。然后每次调用这个对象时，都会改变这个对象中数据成员的状态
func RoundRobinProxySwitcher(proxyURLs ...string) (ProxyFunc, error) {
	if len(proxyURLs) < 1 {
		return nil, errors.New("proxy URL list is empty")
	}
	var urls []*url.URL
	for _, v := range proxyURLs {
		parsedUrl, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		urls = append(urls, parsedUrl)
	}

	rrs := roundRobinSwitcher{
		proxyURls: urls,
		index:     0,
	}
	return (&rrs).GetProxy, nil
}

// GetProxy 作用很像是迭代器
func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	// 这里类似于C语言迭代器执行后置++，返回旧值
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURls[index%uint32(len(r.proxyURls))]
	return u, nil
}
