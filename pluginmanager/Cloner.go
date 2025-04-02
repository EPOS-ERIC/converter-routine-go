package pluginmanager

import (
	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/epos-eu/converter-routine/loggers"
	"gopkg.in/src-d/go-git.v4"
)

func CloneRepository(plugin model.Plugin, options git.CloneOptions) error {
	_, err := git.PlainClone(PluginsPath+plugin.ID, false, &options)
	if err != nil {
		loggers.CRON_LOGGER.Error("Error cloning repository", "plugin", plugin, "error", err)
		return err
	}
	return nil
}
