package proxy

import (
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	proxies commonProxiesBunch
	current int32
}

func NewRoundRobin() *RoundRobin {
	r := &RoundRobin{
		proxies: make(commonProxiesBunch, 0),
		current: -1,
	}

	return r
}

func (r *RoundRobin) Add(url *url.URL) {
	r.proxies = append(r.proxies, NewProxy(url))
	atomic.AddInt32(&r.current, 1)
}

func (r *RoundRobin) Next() (*Proxy, error) {
	next := atomic.AddInt32(&r.current, 1) % int32(len(r.proxies))
	atomic.StoreInt32(&r.current, next)
	return getAvailableProxy(r.proxies, int(next))
}
