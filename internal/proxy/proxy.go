package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

func NewProxy(addr *url.URL) *Proxy {
	return &Proxy{
		proxy: httputil.NewSingleHostReverseProxy(addr),
	}
}

type Proxy struct {
	proxy *httputil.ReverseProxy
	load  int32
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
