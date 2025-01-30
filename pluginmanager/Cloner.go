package pluginmanager

import (
	"log"

	"github.com/epos-eu/converter-routine/dao/model"
	"gopkg.in/src-d/go-git.v4"
)

func CloneRepository(obj model.Softwaresourcecode, options git.CloneOptions) error {
	_, err := git.PlainClone(PluginsPath+obj.InstanceID, false, &options)
	if err != nil {
		log.Printf("Error cloning repository %+v: %v\n", obj, err)
		return err
	}
	return nil
}
