package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epos-eu/converter-routine/connection"
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

// Endpoints
func serviceInit(cs *cronservice.CronService) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Start sync of the plugins endpoint
	r.GET("/sync", func(c *gin.Context) {
		go cs.Task()
		c.JSON(200, "Sync started")
	})

	// Check health (db connection)
	r.GET("/health", healthCheck)

	err := r.Run(":8080")
	panic(err)
}

func healthCheck(c *gin.Context) {
	err := health(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unhealthy: "+err.Error())
		return
	} else {
		c.String(http.StatusOK, "Healthy")
		return
	}
}
func health(c *gin.Context) error {
	// Check the connection to the db
	_, err := connection.Connect()
	if err != nil {
		return fmt.Errorf("can't connect to database")
	}

	return nil
}
