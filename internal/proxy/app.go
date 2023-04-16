package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	loadbalancers "github.com/hiteshrepo/load-balancer/internal/proxy/loadBalancers"
)

type backend struct {
	Url string `json:"url"`
}

type App struct {
	r loadbalancers.RoundRobin
}

func (a *App) Start() {
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

	a.r.Add(v)

	c.String(http.StatusOK, "proxy registered")
}

func (a *App) serveproxy(c *gin.Context) {
	p, err := a.r.Next()
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
		return
	}

	fmt.Printf("current load(%d) \n", p.GetLoad())

	p.ServeHTTP(c.Writer, c.Request)
}
