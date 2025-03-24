package pluginmanager

import (
	"log"

	"github.com/epos-eu/converter-routine/dao/model"
	"gopkg.in/src-d/go-git.v4"
)

func CloneRepository(plugin model.Plugin, options git.CloneOptions) error {
	_, err := git.PlainClone(PluginsPath+plugin.ID, false, &options)
	if err != nil {
		log.Printf("Error cloning repository %+v: %v\n", plugin, err)
		return err
	}
	return nil
}
