package main

import (
	"os"
	"os/signal"
	"syscall"

	rsa "github.com/hiteshrepo/load-balancer/internal/restSimpleApp"
)

func main() {
	app := rsa.App{}

	app.Start()

	<-interrupt()

	app.Stop()

}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}
