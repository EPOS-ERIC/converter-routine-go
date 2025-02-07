package pluginmanager

import (
	"log"
	"os"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/orms"
)

const PluginsPath = "./plugins/"

func Updater() ([]orms.SoftwareSourceCode, error) {
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

func installAndUpdate(sscs []orms.SoftwareSourceCode, branch bool) []orms.SoftwareSourceCode {
	sscs = CloneOrPull(sscs, branch)

	installed := make([]orms.SoftwareSourceCode, 0, len(sscs))
	// for each installed ssc
	for _, ssc := range sscs {
		err := UpdateDependencies(ssc)
		if err != nil {
			// if there is an error getting the dependencies don't consider the plugin as installed
			log.Printf("Error while getting dependencies for %v: %v", ssc.UID, err)
		} else {
			installed = append(installed, ssc)
		}
	}

	return installed
}

type VersionType string

const (
	branch = VersionType("BRANCH")
	tag    = VersionType("TAG")
)
