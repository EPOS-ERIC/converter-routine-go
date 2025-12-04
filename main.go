package main

import (
	"context"

	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/epos-eu/converter-routine/db"
	"github.com/epos-eu/converter-routine/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.Init(); err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// start the cron
	cs := cronservice.NewCronService()
	go cs.Run(ctx)

	// start the service
	go server.serviceInit(cs)

	// block the main goroutine
	select {}
}

type syncHandler struct {
	cs *cronservice.CronService
}
