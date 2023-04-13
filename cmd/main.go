package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hiteshrepo/load-balancer/internal"
)

func main() {
	app := internal.App{}

	app.Start()

	<-interrupt()

	app.Stop()

}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}
