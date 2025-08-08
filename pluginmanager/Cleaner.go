package pluginmanager

import (
	"errors"
	"os"
	"path"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/epos-eu/converter-routine/logging"
	"gorm.io/gorm"
)

var log = logging.Get("plugin_manager")

// CleanPlugins cleans the plugin directory removing all the installations that are not plugins that are currently in the db
func CleanPlugins() error {
	plugins, err := connection.GetPlugins()
	if err != nil {
		return err
	}

	// transform the plugins into a map for fast access
	m := map[string]struct{}{}
	for _, plugin := range plugins {
		m[plugin.ID] = struct{}{}
	}

	// get all the plugin dirs contents
	contents, err := os.ReadDir(PluginsPath)
	if err != nil {
		return err
	}

	for _, file := range contents {
		if _, ok := m[file.Name()]; !ok {
			// make sure that the file is a directory
			if !file.Type().IsDir() {
				log.Error("found an unknown file in the plugin dir", "unknown file name", file.Name())
				continue
			}

			// get the plugin with this id and remove it if it does not exist
			_, err := connection.GetPluginById(file.Name())
			if err != nil {
				// if a plugin with this id does not exist, we can clean it
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// clean the directory removing it

					log.Info("cleaning directory for plugin", "directory name", file.Name())
					err := cleanDir(file.Name())
					if err != nil {
						log.Error("error cleaning plugin", "error", err, "directory name", file.Name())
						continue
					}
				} else { // some other error
					log.Error("error getting pulgin from dir name", "error", err, "directory name", file.Name())
					continue
				}
			}
		}
	}

	return nil
}

// CleanDir deletes a directory and all of its contents from the plugin dir
func cleanDir(name string) error {
	err := os.RemoveAll(path.Join(PluginsPath, name))
	if err != nil {
		return err
	}
	return nil
}

// CleanPlugin cleans a plugin installation and set installed to false for that plugin, then returns it
func CleanPlugin(id string) (plugin model.Plugin, err error) {
	// check that this plugin exists
	plugin, err = connection.GetPluginById(id)
	if err != nil {
		return plugin, err
	}

	// clean its dir
	err = cleanDir(id)
	if err != nil {
		return plugin, err
	}

	// set the installation status to false
	err = connection.SetPluginInstalledStatus(id, false)
	if err != nil {
		return plugin, err
	}

	return plugin, nil
}
