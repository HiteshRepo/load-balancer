package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type backend struct {
	Url string `json:"url"`
}

type App struct {
	proxies commonProxiesBunch
	current int32
}

func (a *App) Start() {
	a.proxies = make(commonProxiesBunch, 0)

	router := gin.Default()
	router.POST("/addurl", a.addproxy)
	router.GET("/", a.serveproxy)

	port := "9091"
	if len(os.Getenv("PROXY_PORT")) > 0 {
		port = os.Getenv("PROXY_PORT")
	}
	router.Run(fmt.Sprintf("localhost:%s", port))
}

func (a *App) Stop() {

}

func (a *App) addproxy(c *gin.Context) {
	var be backend
	if err := c.BindJSON(&be); err != nil {
		return
	}

	v, err := url.Parse(be.Url)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid url(%s) in the request", be.Url))
		return
	}

	p := NewProxy(v)

	a.proxies = append(a.proxies, p)
	atomic.AddInt32(&a.current, 1)

	c.String(http.StatusOK, "proxy registered")
}

func (a *App) serveproxy(c *gin.Context) {
	p, err := getAvailableProxy(a.proxies, 0)
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
		return
	}

	fmt.Printf("current load(%d) \n", p.GetLoad())

	p.ServeHTTP(c.Writer, c.Request)
}
