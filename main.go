package main

import (
	"context"
	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/gin-gonic/gin"
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

// get endpoint to start syncing
func serviceInit(cs *cronservice.CronService) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/sync", func(c *gin.Context) {
		go cs.Task()
		c.JSON(200, "Sync started")
	})
	err := r.Run(":8080")
	panic(err)
}
