package cronservice

import (
	"context"
	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/pluginmanager"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
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
	log.Printf("Cron task started at %v\n", time.Now())
	installedRepos := pluginmanager.Updater()

	newPlugins, err := connection.GeneratePlugins(installedRepos)
	if err != nil {
		log.Printf("Error generating new plugins: %v\n", err)
		return
	}

	log.Printf("Found %d new plugins\n", len(newPlugins))

	if len(newPlugins) > 0 {
		if err := connection.InsertPlugins(newPlugins); err != nil {
			log.Printf("Error inserting new plugins: %v\n", err)
		} else {
			log.Println("Successfully inserted new plugins")
		}
	} else {
		log.Println("Plugins up to date")
	}

	newPluginsRelations, err := connection.GeneratePluginsRelations()
	if err != nil {
		log.Printf("Error generating new plugin relations: %v\n", err)
		return
	}

	log.Printf("Found %d new plugin relations\n", len(newPluginsRelations))

	if len(newPluginsRelations) > 0 {
		if err := connection.InsertPluginsRelations(newPluginsRelations); err != nil {
			log.Printf("Error inserting new plugin relations: %v\n", err)
		} else {
			log.Println("Successfully inserted new plugin relations")
		}
	} else {
		log.Println("Plugin relations up to date")
	}

	log.Printf("Cron task ended at %v\n", time.Now())
	taskMutex.Unlock()
}
