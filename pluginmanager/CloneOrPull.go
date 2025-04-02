package pluginmanager

import (
	"fmt"
	"os"

	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/epos-eu/converter-routine/loggers"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CloneOrPull a plugin
func CloneOrPull(plugin model.Plugin) error {
	// Determine the reference name based on the provided options
	var referenceName plumbing.ReferenceName
	if plugin.VersionType == "branch" {
		referenceName = plumbing.NewBranchReferenceName(plugin.Version)
	} else {
		referenceName = plumbing.NewTagReferenceName(plugin.Version)
	}

	// Define clone and pull options
	cloneOptions := git.CloneOptions{
		URL:           plugin.Repository,
		ReferenceName: referenceName,
	}
	pullOptions := git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: referenceName,
		SingleBranch:  true,
	}

	// Construct the repository path using the instance ID
	repoPath := PluginsPath + plugin.ID

	// Check if the repository directory exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		loggers.CRON_LOGGER.Info("Repository does not exist, cloning it", "repository", plugin.Repository)
		// If the repository does not exist, clone it
		err = CloneRepository(plugin, cloneOptions)
		if err != nil {
			return fmt.Errorf("error while cloning %v: %v", plugin.ID, err)
		}
	} else {
		// Define checkout options
		checkoutOptions := git.CheckoutOptions{
			Branch: referenceName,
		}

		// Checkout the specified branch
		if err := Checkout(plugin, checkoutOptions); err != nil {
			return fmt.Errorf("error checking out branch %v: %v", referenceName, err)
		}

		loggers.CRON_LOGGER.Info("Repository exists, pulling it", "repository", plugin.Repository)
		// If the repository exists, attempt to pull the latest changes
		if err := PullRepository(plugin, pullOptions); err != nil {
			return fmt.Errorf("error pulling: %v", err)
		}
	}

	return nil
}
