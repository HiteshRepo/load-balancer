package types

import (
	"fmt"
)

type proxiesBunch interface {
	Len() int
	Get(idx int) *Proxy
}

func GetAvailableProxy(proxies proxiesBunch, marker int) (*Proxy, error) {
	for i := 0; i < proxies.Len(); i++ {
		tryProxy := (marker + i) % proxies.Len()
		p := proxies.Get(tryProxy)
		if p != nil && p.IsAvailable() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}

type CommonProxiesBunch []*Proxy

func (b CommonProxiesBunch) Len() int           { return len(b) }
func (b CommonProxiesBunch) Get(idx int) *Proxy { return b[idx] }
