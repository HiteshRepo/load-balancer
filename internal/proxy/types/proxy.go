package types

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/hiteshrepo/load-balancer/internal/proxy/health"
)

type Proxy struct {
	proxy  *httputil.ReverseProxy
	health *health.ProxyHealth
	load   int32
}

func NewProxy(addr *url.URL) *Proxy {
	return &Proxy{
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		health: health.NewProxyHealth(addr),
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)
	time.Sleep(10 * time.Second)
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}

func (p *Proxy) IsAvailable() bool {
	return p.health.IsAvailable()
}

func (p *Proxy) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	p.health.SetHealthCheck(check, period)
}
