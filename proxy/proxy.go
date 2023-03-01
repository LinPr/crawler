package proxy

import (
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

type RoundRobinBalancer struct {
	lock      sync.Mutex
	proxyURls []*url.URL
	index     uint32
}

// 这是一个闭包，
// 在go语言中函数作为一等公民，函数可以当做任何参数或者返回值传来传去（类似C++中的std::function，或者仿函数对象）
// 使用方法为，根据这个闭包生成一个对象。然后每次调用这个对象时，都会改变这个对象中数据成员的状态
func NewRoundRobinBalancer(proxyURLs ...string) (ProxyFunc, error) {
	if len(proxyURLs) == 0 {
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

	rrb := RoundRobinBalancer{
		proxyURls: urls,
		index:     0,
	}
	return (&rrb).GetNextServer, nil
}

// GetProxy 作用很像是迭代器,类似C语言的后置++
func (r *RoundRobinBalancer) GetNextServer(pr *http.Request) (*url.URL, error) {
	r.lock.Lock()
	u := r.proxyURls[r.index]
	r.index = (r.index + 1) % uint32(len(r.proxyURls))
	r.lock.Unlock()
	return u, nil
}

type ConsistentHashBalancer struct {
	lock      sync.RWMutex
	replicas  int
	proxyURLs []*url.URL
	circle    map[uint32]*url.URL
}

func NewConsistentHashBalancer(replicas int, proxyURLs ...string) (ProxyFunc, error) {
	if len(proxyURLs) == 0 {
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

	chb := ConsistentHashBalancer{
		replicas:  replicas,
		proxyURLs: urls,
		circle:    make(map[uint32]*url.URL),
	}
	// 每个proxy都在哈希环上创建指定的副本数量
	for _, v := range urls {
		for i := 0; i < replicas; i++ {
			hash := crc32.ChecksumIEEE([]byte(v.String() + strconv.Itoa(i)))
			chb.circle[hash] = v
		}
	}
	return (&chb).GetNextServer, nil

}

// 读取哈希环上下一个可用的虚拟 proxy 节点，并通过map映射到真实的 proxy
func (c *ConsistentHashBalancer) GetNextServer(pr *http.Request) (*url.URL, error) {
	// 只读操作加读锁
	c.lock.RLock()
	defer c.lock.RUnlock()

	//这里因为我只有一台机器，多代理，是一对多的无状态模型，以每次请求访问的时间作为key，访问哈希环,
	//如果是多对多，且有状态可用机器 IP + port 作key
	key, _ := time.Now().MarshalBinary()
	hash := crc32.ChecksumIEEE(key)
	for k := range c.circle {
		if k >= hash {
			return c.circle[k], nil
		}
	}

	// 若这个哈希值比所有的虚拟节点都大, 则找到key值最小的节点
	var minKey uint32 = ^uint32(0)
	for k, _ := range c.circle {
		if k < minKey {
			minKey = k
		}
	}
	return c.circle[minKey], nil
}

func (c *ConsistentHashBalancer) RemoveServer(proxyURL string) error {
	// 删除操作加写锁
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := crc32.ChecksumIEEE([]byte(proxyURL + strconv.Itoa(i)))
		delete(c.circle, hash)
	}

	parsedUrl, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("url.Parse() error: %v\n", err)
		return err
	}
	for k, v := range c.proxyURLs {
		if v == parsedUrl {
			c.proxyURLs = append(c.proxyURLs[:k], c.proxyURLs[k+1:]...)
			break
		}
	}
	return nil
}

func (c *ConsistentHashBalancer) AddServer(proxyURL string) error {
	// 增加操作加写锁
	c.lock.Lock()
	defer c.lock.Unlock()

	parsedUrl, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("url.Parse() error: %v\n", err)
		return err
	}
	c.proxyURLs = append(c.proxyURLs, parsedUrl)
	for i := 0; i < c.replicas; i++ {
		hash := crc32.ChecksumIEEE([]byte(proxyURL + strconv.Itoa(i)))
		c.circle[hash] = parsedUrl
	}

	return nil
}
