package pluginmanager

import (
	"fmt"
	"log"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/dao/model"
)

const PluginsPath = "./plugins/"

func Updater() error {
	plugins, err := connection.GetPlugins()
	if err != nil {
		return err
	}
	if len(plugins) <= 0 {
		return fmt.Errorf("no plugins found while updating")
	}

	log.Printf("Found %d plugins\n", len(plugins))

	for _, plugin := range plugins {
		err := installAndUpdate(plugin)
		if err != nil {
			log.Printf("error while installing and updating plugin %s: %s", plugin.ID, err)

			// if there has been an error, don't consider this plugin as installed
			connection.SetPluginInstalledStatus(plugin.ID, false)
			continue
		}

		connection.SetPluginInstalledStatus(plugin.ID, true)
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
