package cronservice

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Fatal(err)
	}
	ds.cron.Start()
	<-ctx.Done()
	cronTask := ds.cron.Stop()
	<-cronTask.Done()
}

var taskMutex sync.Mutex

// periodic Task
func (ds *CronService) Task() {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	log.Printf("Cron task started at %v\n", time.Now())

	err := pluginmanager.Updater()
	if err != nil {
		log.Printf("Error updating plugins: %v\n", err)
		return
	}

	log.Printf("Cron task ended at %v\n", time.Now())
}
