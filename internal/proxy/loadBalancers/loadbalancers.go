package loadBalancers

import (
	"net/url"

	"github.com/hiteshrepo/load-balancer/internal/proxy/types"
)

type LoadBalancer interface {
	Next() (*types.Proxy, error)
	Add(url *url.URL, weight int32)
}
