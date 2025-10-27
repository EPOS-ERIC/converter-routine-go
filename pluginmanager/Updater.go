package pluginmanager

import (
	"fmt"

	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/epos-eu/converter-routine/db"
)

const PluginsPath = "./plugins/"

func SyncPlugins() error {
	plugins, err := db.GetPlugins()
	if err != nil {
		return err
	}
	if len(plugins) <= 0 {
		log.Warn("No plugins found while updating")
		return nil
	}

	log.Info("Found plugins", "count", len(plugins))

	for _, plugin := range plugins {
		err := SyncPlugin(plugin)
		// if a sync fails don't fail the whole task
		if err != nil {
			continue
		}
	}
	return nil
}

func SyncPlugin(plugin model.Plugin) error {
	err := installAndUpdate(plugin)
	if err != nil {
		log.Error("Error while installing and updating plugin", "pluginID", plugin.ID, "error", err)
		// if there has been an error, don't consider this plugin as installed
		newErr := db.SetPluginInstalledStatus(plugin.ID, false)
		if newErr != nil {
			return newErr
		}
		return err
	}

	err = db.SetPluginInstalledStatus(plugin.ID, true)
	if err != nil {
		return err
	}
	return nil
}

func installAndUpdate(plugin model.Plugin) error {
	err := CloneOrPull(plugin)
	if err != nil {
		return err
	}

	err = UpdateDependencies(plugin)
	if err != nil {
		// if there is an error getting the dependencies don't consider the plugin as installed
		return fmt.Errorf("error while updating dependencies for %v: %v", plugin.ID, err)
	}

	return nil
}
