package loadBalancers

import (
	"net/url"
	"sync/atomic"

	"github.com/hiteshrepo/load-balancer/internal/proxy/types"
)

type RoundRobin struct {
	proxies types.CommonProxiesBunch
	current int32
}

func NewRoundRobin() *RoundRobin {
	r := &RoundRobin{
		proxies: make(types.CommonProxiesBunch, 0),
		current: -1,
	}

	return r
}

func (r *RoundRobin) Add(url *url.URL) {
	r.proxies = append(r.proxies, types.NewProxy(url))
	atomic.AddInt32(&r.current, 1)
}

func (r *RoundRobin) Next() (*types.Proxy, error) {
	next := atomic.AddInt32(&r.current, 1) % int32(len(r.proxies))
	atomic.StoreInt32(&r.current, next)
	return types.GetAvailableProxy(r.proxies, int(next))
}
