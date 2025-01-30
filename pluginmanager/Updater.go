package pluginmanager

import (
	"log"
	"os"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/dao/model"
)

const PluginsPath = "./plugins/"

func Updater() ([]model.Softwaresourcecode, error) {
	scss, err := connection.GetSoftwareSourceCodes()
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d software source codes\n", len(scss))

	// get the type of the version from the env variables, if not set or set wrong treat the version as branch
	versionType := os.Getenv("PLUGINS_VERSION_TYPE")

	switch VersionType(versionType) {
	case tag:
		return installAndUpdate(scss, false), nil
	default: // branch
		return installAndUpdate(scss, true), nil
	}
}

func installAndUpdate(sscs []model.Softwaresourcecode, branch bool) []model.Softwaresourcecode {
	sscs = CloneOrPull(sscs, branch)

	// for each installed ssc
	for i, ssc := range sscs {
		err := UpdateDependencies(ssc)
		if err != nil {
			// if there is an error getting the dependencies don't consider the plugin as installed
			sscs = append(sscs[:i], sscs[i+1:]...)
			log.Printf("Error while getting dependencies for %v: %v", ssc.UID, err)
		}
	}

	return sscs
}

type VersionType string

const (
	branch = VersionType("BRANCH")
	tag    = VersionType("TAG")
)
