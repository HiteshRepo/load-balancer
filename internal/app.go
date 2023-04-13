package internal

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func (a *App) Start() {
	router := gin.Default()
	router.GET("/", response)

	port := "9090"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	router.Run(fmt.Sprintf("localhost:%s", port))
}

func (a *App) Stop() {

}

func response(c *gin.Context) {
	var podIp string
	if len(os.Getenv("POD_NAME")) > 0 {
		podIp = os.Getenv("POD_NAME")
	}
	c.String(http.StatusOK, "responding from a simple rest app with pod name: %s", podIp)
}
