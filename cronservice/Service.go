package cronservice

import (
	"context"
	"os"
	"sync"

	"github.com/epos-eu/converter-routine/loggers"
	"github.com/epos-eu/converter-routine/pluginmanager"
	"github.com/robfig/cron/v3"
)

// cronRunner interface for cron
type cronRunner interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Start()
	Stop() context.Context
}

// CronService cron service struct
type CronService struct {
	cron cronRunner
}

// NewCronService returns new cron service
func NewCronService() *CronService {
	return &CronService{
		cron: cron.New(),
	}
}

const (
	timePatterns = "*/5 * * * *"
)

// Run starts service
func (ds *CronService) Run(ctx context.Context) {
	// Execute the task immediately
	ds.Task()

	// Schedule the task to run every 5 minutes
	if _, err := ds.cron.AddFunc(timePatterns, ds.Task); err != nil {
		loggers.CRON_LOGGER.Error("Failed to schedule cron task", "error", err)
		os.Exit(1)
	}
	ds.cron.Start()
	<-ctx.Done()
	cronTask := ds.cron.Stop()
	<-cronTask.Done()
}

var taskMutex sync.Mutex

// Task is the periodic cron task
func (ds *CronService) Task() {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	loggers.CRON_LOGGER.Info("Cron task started")

	err := pluginmanager.Updater()
	if err != nil {
		loggers.CRON_LOGGER.Error("Error updating plugins", "error", err)
		return
	}

	loggers.CRON_LOGGER.Info("Cron task ended")
}
