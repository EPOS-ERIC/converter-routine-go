package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/epos-eu/converter-routine/pluginmanager"
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

	// Delete plugin directory
	r.GET("/clean/:id", cleanPlugin)

	err := r.Run(":8080")
	panic(err)
}

func cleanPlugin(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.String(http.StatusBadRequest, "Missing parameter 'id'")
		return
	}
	plugin, err := connection.GetPluginById(id)
	if err != nil {
		c.String(http.StatusBadRequest, "No plugin found with id: %v", id)
		return
	}
	err = os.RemoveAll(path.Join(pluginmanager.PluginsPath, plugin.ID))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error cleaning plugin: %v", plugin)
		return
	}
	c.String(http.StatusOK, "Plugin cleaned")
}

func healthCheck(c *gin.Context) {
	err := health()
	if err != nil {
		c.String(http.StatusInternalServerError, "Unhealthy: "+err.Error())
		return
	} else {
		c.String(http.StatusOK, "Healthy")
		return
	}
}

func health() error {
	// Check the connection to the db
	_, err := connection.ConnectMetadata()
	if err != nil {
		return fmt.Errorf("can't connect to Metadata database")
	}

	_, err = connection.ConnectConverter()
	if err != nil {
		return fmt.Errorf("can't connect to Converter database")
	}

	return nil
}
