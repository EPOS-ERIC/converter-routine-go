package pluginmanager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/epos-eu/converter-routine/connection"
	"github.com/epos-eu/converter-routine/orms"
)

// UpdateDependencies installs (or updates) the dependencies for a plugin depending on its runtime.
func UpdateDependencies(ssc orms.SoftwareSourceCode) error {
	lang, err := connection.GetSoftwareSourceCodeProgrammingLanguage(ssc.Instance_id)
	if err != nil {
		return fmt.Errorf("error getting programming language for SoftwareSourceCode with instanceId %s: %w", ssc.Instance_id, err)
	}

	log.Printf("Updating dependencies for SoftwareSourceCode %s...", ssc.Instance_id)

	switch lang {
	case "Go", "Java":
		// no dependencies handling for java and go plugins
		log.Printf("\tDONE: No dependencies to update")
		return nil
	case "Python":
		return handlePyhonDependencies(ssc)
	default:
		return fmt.Errorf("error: unknown runtime: %s", lang)
	}
}

// handlePyhonDependencies sets up a Venv python environment and then installs the dependencies
func handlePyhonDependencies(ssc orms.SoftwareSourceCode) error {
	path := filepath.Join(PluginsPath, ssc.Instance_id)

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

	log.Println("\tPython venv set up correctly")

	// Execute the shell command
	cmd = exec.Command("venv/bin/pip", "install", "-r", "requirements.txt")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		log.Fatalf("\tError installing dependencies: %v", err)
	}

	log.Println("\tPython dependencies installed successfully")

	return nil
}
