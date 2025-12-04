package pluginmanager

import (
	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/go-git/go-git/v5"
)

func CloneRepository(plugin model.Plugin, options git.CloneOptions) error {
	_, err := git.PlainClone(PluginsPath+plugin.ID, false, &options)
	if err != nil {
		log.Error("Error cloning repository", "plugin", plugin, "error", err)
		return err
	}
	return nil
}
