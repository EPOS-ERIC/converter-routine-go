package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/epos-eu/converter-routine/cronservice"
	"github.com/epos-eu/converter-routine/db"
	"github.com/epos-eu/converter-routine/logging"
	"github.com/epos-eu/converter-routine/pluginmanager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var log = logging.Get("api")

//go:embed openapi.json
var openAPISpec []byte

// slogGinMiddleware creates a structured logging middleware using the logging package
func slogGinMiddleware() gin.HandlerFunc {
	// Get a logger specifically for HTTP requests
	httpLog := logging.Get("http")

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if strings.Contains(path, "actuator/health") {
			return
		}

		latency := time.Since(start)
		clientIP := c.ClientIP()
		if raw != "" {
			path = path + "?" + raw
		}

		status := c.Writer.Status()
		var level slog.Level
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		default:
			level = slog.LevelInfo
		}

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.String("client_ip", clientIP),
			slog.Duration("latency", latency),
			slog.Int64("response_size", int64(c.Writer.Size())),
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			attrs = append(attrs,
				slog.String("error", err.Error()),
				slog.Any("error_type", err.Type),
			)
		}

		httpLog.LogAttrs(c.Request.Context(), level, "HTTP request", attrs...)
	}
}

// customRecoveryMiddleware handles panics with structured logging
func customRecoveryMiddleware() gin.HandlerFunc {
	recoveryLog := logging.Get("recovery")

	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		recoveryLog.Error("panic recovered",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
			slog.Any("panic", recovered),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// Endpoints
func serviceInit(cs *cronservice.CronService) {
	//	@title		Converter Routine API
	//	@version	1.0
	//	@BasePath	/api/converter-routine/v1

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()

	r.Use(customRecoveryMiddleware())

	r.Use(slogGinMiddleware())

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

	plugin, err := db.GetPluginById(id)
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
		log.Error("health check failed", "error", err)
		c.String(http.StatusServiceUnavailable, "Unhealthy")
		return
	} else {
		c.String(http.StatusOK, "Healthy")
		return
	}
}

func health() error {
	db := db.Get()

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("can't get underlying sql.DB: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("can't ping database: %w", err)
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
