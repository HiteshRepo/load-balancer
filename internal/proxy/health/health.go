package health

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"
)

const (
	defaultHealthCheckTimeout = 10 * time.Second
	defaultHealthCheckPeriod  = 10 * time.Second
)

func NewProxyHealth(origin *url.URL) *ProxyHealth {
	h := &ProxyHealth{
		origin:      origin,
		check:       defaultHealthCheck,
		period:      defaultHealthCheckPeriod,
		cancel:      make(chan struct{}),
		isAvailable: defaultHealthCheck(origin),
	}
	h.run()

	return h
}

type ProxyHealth struct {
	origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
}

func (h *ProxyHealth) IsAvailable() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.isAvailable
}

func (h *ProxyHealth) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.stop()
	h.check = check
	h.period = period
	h.cancel = make(chan struct{})
	h.isAvailable = h.check(h.origin)
	h.run()
}

func (h *ProxyHealth) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stop()
}

func (h *ProxyHealth) run() {
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		isAvailable := h.check(h.origin)
		h.isAvailable = isAvailable
	}

	go func() {
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-h.cancel:
				t.Stop()
				return
			}
		}
	}()
}

func (h *ProxyHealth) stop() {
	if h.cancel != nil {
		h.cancel <- struct{}{}
		close(h.cancel)
		h.cancel = nil
	}
}

var defaultHealthCheck = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, defaultHealthCheckTimeout)
	if err != nil {
		fmt.Printf("%s -> url(%s) is unavailable.\n", time.Now().Format("2017-09-07 17:06:06"), addr.String())
		return false
	}

	fmt.Printf("%s -> url(%s) is available.\n", time.Now().Format("2017-09-07 17:06:06"), addr.String())

	_ = conn.Close()
	return true
}
