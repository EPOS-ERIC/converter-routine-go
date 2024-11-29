package pluginmanager

import (
	"github.com/epos-eu/converter-routine/orms"
	"gopkg.in/src-d/go-git.v4"
	"log"
)

func CloneRepository(obj orms.SoftwareSourceCode, options git.CloneOptions) error {
	log.Println(obj.GetRuntimeplatform())

	_, err := git.PlainClone(PluginsPath+obj.GetInstance_id(), false, &options)
	if err != nil {
		log.Printf("Error cloning repository %v: %v\n", obj.Uid, err)
		return err
	}
	return nil
}
