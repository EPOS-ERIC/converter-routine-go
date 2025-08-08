package main

import (
	"context"

	"github.com/epos-eu/converter-routine/cronservice"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the cron
	cs := cronservice.NewCronService()
	go cs.Run(ctx)

	// start the service
	go serviceInit(cs)

	// block the main goroutine
	select {}
}

type syncHandler struct {
	cs *cronservice.CronService
}
