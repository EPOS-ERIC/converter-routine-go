package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/epos-eu/converter-routine/loggers"
	"github.com/epos-eu/converter-routine/pluginmanager"
	"github.com/gin-gonic/gin"
)

func main() {
	// init json logging
	loggers.InitSlog()

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

// Embed the OpenAPI 3.0 JSON specification file
//
//go:embed openapi.json
var openAPISpec []byte

// Endpoints
func serviceInit(cs *cronservice.CronService) {
	//	@title		Converter Routine API
	//	@version	1.0
	//	@BasePath	/api/converter-routine/v1

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	syncHandler := syncHandler{cs}
	v1 := r.Group("/api/converter-routine/v1")
	{
		v1.GET("/sync", syncHandler.sync)

		// Check health (db connection)
		v1.GET("/health", healthCheck)

		// Delete plugin directory
		v1.GET("/clean/:plugin_id", cleanPlugin)

		// Swagger json
		v1.GET("/api-docs", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/json", openAPISpec)
		})

		// TODO: reinstall endpoint
	}

	err := r.Run(":8080")
	panic(err)
}

type syncHandler struct {
	cs *cronservice.CronService
}

// Start a sync
//
//	@Summary		Start sync process
//	@Description	Initiates a background synchronization task and returns immediately.
//	@Tags			Converter Routine
//	@Produce		json
//	@Success		200	{string}	string	"Sync started"
//	@Router			/sync [get]
func (s *syncHandler) sync(c *gin.Context) {
	go s.cs.Task()
	c.JSON(200, "Sync started")
}

// cleanPlugin cleans the plugin directory in the volume
//
//	@Summary		Cleans the plugin installation
//	@Description	Removes the plugin directory for the specified plugin ID.
//	@Tags			Converter Routine
//	@Produce		json
//	@Param			plugin_id	path		string		true	"Plugin ID"
//	@Success		200	{string}	string		"Plugin cleaned"
//	@Failure		400	{object}	HTTPError	"Missing parameter 'plugin_id'"
//	@Failure		404	{object}	HTTPError	"Plugin not found"
//	@Failure		500	{object}	HTTPError	"Error cleaning plugin"
//	@Router			/clean/{plugin_id} 	[get]
func cleanPlugin(c *gin.Context) {
	id, ok := c.Params.Get("plugin_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'plugin_id'"})
		return
	}
	plugin, err := connection.GetPluginById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No plugin found with plugin_id: " + id})
		return
	}
	err = os.RemoveAll(path.Join(pluginmanager.PluginsPath, plugin.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cleaning plugin: " + err.Error()})
		return
	}
	err = connection.SetPluginInstalledStatus(id, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cleaning plugin: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Plugin cleaned")
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

// HTTPError is used just by swag
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
