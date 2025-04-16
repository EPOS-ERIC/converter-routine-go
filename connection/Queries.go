package connection

import (
	"github.com/epos-eu/converter-routine/dao/model"
)

func GetPlugins() ([]model.Plugin, error) {
	db, err := ConnectConverter()
	if err != nil {
		return nil, err
	}

	var listOfPlugins []model.Plugin
	err = db.Find(&listOfPlugins).Error
	if err != nil {
		return nil, err
	}
	return listOfPlugins, nil
}

func GetPluginById(pluginId string) (model.Plugin, error) {
	var plugin model.Plugin
	db, err := ConnectConverter()
	if err != nil {
		return plugin, err
	}
	err = db.Model(&plugin).Where("id = ?", pluginId).First(&plugin).Error
	if err != nil {
		return plugin, err
	}
	return plugin, nil
}

func GetPluginRelations() ([]model.PluginRelation, error) {
	db, err := ConnectConverter()
	if err != nil {
		return nil, err
	}

	var listOfPluginRelations []model.PluginRelation
	err = db.Find(&listOfPluginRelations).Error
	if err != nil {
		return nil, err
	}
	return listOfPluginRelations, nil
}

func SetPlugins(ph []model.Plugin) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	// Truncate plugins table
	err = db.Exec("TRUNCATE plugin CASCADE").Error
	if err != nil {
		return err
	}

	// Insert new plugins
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}
	return nil
}

func SetPluginsRelations(ph []model.PluginRelation) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	// Truncate plugin_relations table
	err = db.Exec("TRUNCATE plugin_relations CASCADE").Error
	if err != nil {
		return err
	}

	// Insert new plugin relations
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertPlugins(plugins []model.Plugin) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	err = db.Create(&plugins).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertPluginsRelations(pluginRelations []model.PluginRelation) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	err = db.Create(&pluginRelations).Error
	if err != nil {
		return err
	}
	return nil
}

// SetPluginStatus given an id and a status sets the current installed status of a plugin
func SetPluginInstalledStatus(id string, installed bool) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	// Find the existing plugin record by ID
	var existing model.Plugin
	err = db.First(&existing, "id = ?", id).Error
	if err != nil {
		return err
	}

	existing.Installed = installed

	// Update the existing plugin record with the new data
	err = db.Model(&existing).Select("*").Updates(existing).Error
	if err != nil {
		return err
	}

	return nil
}
