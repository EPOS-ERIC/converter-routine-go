package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/epos-eu/converter-routine/loggers"
	"github.com/epos-eu/converter-routine/pluginmanager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type LogEntry struct {
	Component string `json:"component"`
	Time      string `json:"time"`
	Level     string `json:"level"`
	Status    int    `json:"status"`
	Latency   string `json:"latency"`
	ClientIP  string `json:"client_ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	BodySize  int    `json:"body_size"`
	Error     string `json:"error,omitempty"`
}

// Endpoints
func serviceInit(cs *cronservice.CronService) {
	//	@title		Converter Routine API
	//	@version	1.0
	//	@BasePath	/api/converter-routine/v1

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	r := gin.New()
	r.Use(gin.Recovery())
	// use a custom json logger for gin
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		entry := LogEntry{
			Component: "API",
			Time:      param.TimeStamp.Format(time.RFC3339),
			Level:     "INFO",
			Status:    param.StatusCode,
			Latency:   param.Latency.String(),
			ClientIP:  param.ClientIP,
			Method:    param.Method,
			Path:      param.Path,
			BodySize:  param.BodySize,
			Error:     param.ErrorMessage,
		}

		jsonLog, _ := json.Marshal(entry)
		return string(jsonLog) + "\n"
	}))

	syncHandler := syncHandler{cs}
	v1 := r.Group("/api/converter-routine/v1")
	{
		v1.POST("/sync", syncHandler.sync)
		v1.POST("/sync/:plugin_id", syncHandler.syncPlugin)

		// Check health (db connection)
		v1.GET("/actuator/health", healthCheck)

		// Delete plugin directory
		v1.POST("/clean/:plugin_id", cleanPlugin)

		// reinstall endpoint
		v1.POST("/reinstall/:plugin_id", reinstallPlugin)

		// Swagger json
		v1.GET("/api-docs", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/json", openAPISpec)
		})
	}

	err := r.Run(":8080")
	panic(err)
}

type syncHandler struct {
	cs *cronservice.CronService
}

// Start a sync
//
//	@Summary		Start the async process for updating all plugin
//	@Description	Initiates a background synchronization task and returns immediately.
//	@Tags			Converter Routine
//	@Produce		json
//	@Success		202	{object}	OK
//	@Router			/sync [post]
func (s *syncHandler) sync(c *gin.Context) {
	go s.cs.Task()
	c.JSON(http.StatusAccepted, "Sync started")
}

// Start a sync
//
//	@Summary		Start the sync process for a specific plugin. This request will answer when the plugin sync has either completed or failed.
//	@Description	Initiates a background synchronization task and returns immediately.
//	@Tags			Converter Routine
//	@Param			plugin_id	path	string	true	"Plugin ID"
//	@Produce		json
//	@Success		200	{object}	OK
//	@Router			/sync/{plugin_id} [post]
func (s *syncHandler) syncPlugin(c *gin.Context) {
	id, ok := c.Params.Get("plugin_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'plugin_id'"})
		return
	}

	plugin, err := connection.GetPluginById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found in DB: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync plugin: " + err.Error()})
	}
	err = pluginmanager.SyncPlugin(plugin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin " + id + "synced successfully"})
}

// cleanPlugin cleans the plugin directory in the volume
//
//	@Summary		Cleans the plugin installation
//	@Description	Removes the plugin directory for the specified plugin ID.
//	@Tags			Converter Routine
//	@Produce		json
//	@Param			plugin_id			path		string	true	"Plugin ID"
//	@Success		200					{object}	OK
//	@Failure		400					{object}	HTTPError	"Missing parameter 'plugin_id'"
//	@Failure		404					{object}	HTTPError	"Plugin not found"
//	@Failure		500					{object}	HTTPError	"Error cleaning plugin"
//	@Router			/clean/{plugin_id} 	[post]
func cleanPlugin(c *gin.Context) {
	id, ok := c.Params.Get("plugin_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'plugin_id'"})
		return
	}
	// delete the installation dir
	_, err := pluginmanager.CleanPlugin(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cleaning plugin: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin cleaned"})
}

// reinstallPlugin cleans the plugin directory and syncs the plugin
//
//	@Summary		Cleans the plugin installation and syncs the plugin
//	@Description	Removes the plugin directory for the specified plugin ID and syncs the plugin
//	@Tags			Converter Routine
//	@Produce		json
//	@Param			plugin_id				path		string	true	"Plugin ID"
//	@Success		200						{object}	OK
//	@Failure		400						{object}	HTTPError	"Missing parameter 'plugin_id'"
//	@Failure		404						{object}	HTTPError	"Plugin not found"
//	@Failure		500						{object}	HTTPError	"Error reinstalling plugin"
//	@Router			/reinstall/{plugin_id} 	[post]
func reinstallPlugin(c *gin.Context) {
	id, ok := c.Params.Get("plugin_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'plugin_id'"})
		return
	}

	// clean the plugin
	plugin, err := pluginmanager.CleanPlugin(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cleaning plugin: " + err.Error()})
		return
	}
	// re install it
	err = pluginmanager.SyncPlugin(plugin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error syncing plugin: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin cleaned and reinstalled successfully"})
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
	_, err := connection.ConnectConverter()
	if err != nil {
		return fmt.Errorf("can't connect to Converter database")
	}

	return nil
}

// HTTPError is used just by swag
type HTTPError struct {
	Error string `json:"error" example:"can't connect to database"`
}

// just used by swag
type OK struct {
	Message string `json:"message" example:"sync started"`
}
