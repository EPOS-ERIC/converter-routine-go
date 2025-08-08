package pluginmanager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/epos-eu/converter-routine/dao/model"
)

// UpdateDependencies installs (or updates) the dependencies for a plugin depending on its runtime
func UpdateDependencies(plugin model.Plugin) error {
	switch plugin.Runtime {
	case "binary", "java":
		// No dependencies handling for java and go plugins
		log.Debug("No dependencies to update", "plugin", plugin.ID)
		return nil
	case "python":
		return handlePythonDependencies(plugin)
	default:
		return fmt.Errorf("error: unknown runtime: %s", plugin.Runtime)
	}
}

// handlePythonDependencies sets up a Python venv environment and then installs the dependencies
func handlePythonDependencies(plugin model.Plugin) error {
	path := filepath.Join(PluginsPath, plugin.ID)

	_, err := os.Stat(filepath.Join(path, "requirements.txt"))
	if os.IsNotExist(err) {
		return fmt.Errorf("error installing dependencies: file 'requirements.txt' not found")
	} else if err != nil {
		return fmt.Errorf("error installing dependencies: error while checking existence of 'requirements.txt': %w", err)
	}

	// Initialize the venv environment
	cmd := exec.Command("python3", "-m", "venv", "venv")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating venv environment: %w", err)
	}

	log.Info("Python venv set up correctly", "plugin", plugin.ID)

	// Execute the command to install dependencies
	cmd = exec.Command("venv/bin/pip", "install", "-r", "requirements.txt")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error installing dependencies for plugin %s: %w", plugin.ID, err)
	}

	log.Info("Python dependencies installed successfully", "plugin", plugin.ID)
	return nil
}
