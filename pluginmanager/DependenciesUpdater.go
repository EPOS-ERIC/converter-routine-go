package pluginmanager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/epos-eu/converter-routine/dao/model"
)

// UpdateDependencies installs (or updates) the dependencies for a plugin depending on its runtime.
func UpdateDependencies(plugin model.Plugin) error {
	switch plugin.Runtime {
	case "binary", "java":
		// no dependencies handling for java and go plugins
		// log.Printf("\tDONE: No dependencies to update")
		return nil
	case "python":
		return handlePyhonDependencies(plugin)
	default:
		return fmt.Errorf("error: unknown runtime: %s", plugin.Runtime)
	}
}

// handlePyhonDependencies sets up a Venv python environment and then installs the dependencies
func handlePyhonDependencies(plugin model.Plugin) error {
	path := filepath.Join(PluginsPath, plugin.ID)

	_, err := os.Stat(filepath.Join(path, "requirements.txt"))
	if os.IsNotExist(err) {
		return fmt.Errorf("error installing dependencies: file 'reqirements.txt' not found")
	} else if err != nil {
		return fmt.Errorf("error installing dependencies: error while cheking existance of 'requirements.txt': %w", err)
	}

	// initialize the venv environment
	cmd := exec.Command("python3", "-m", "venv", "venv")
	// set the directory where to execute the command ./plugin/{ssc.Instance_id}
	cmd.Dir = path
	// create the venv environment. if it already exists, nothing will happen
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating venv environment: %w", err)
	}

	// log.Println("\tPython venv set up correctly")

	// Execute the shell command
	cmd = exec.Command("venv/bin/pip", "install", "-r", "requirements.txt")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		log.Fatalf("\tError installing dependencies: %v", err)
	}

	// log.Println("\tPython dependencies installed successfully")

	return nil
}
