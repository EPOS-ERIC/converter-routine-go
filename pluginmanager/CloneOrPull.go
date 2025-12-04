package pluginmanager

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneOrPull(plugin model.Plugin) error {
	var referenceName plumbing.ReferenceName
	switch plugin.VersionType {
	case "branch":
		referenceName = plumbing.NewBranchReferenceName(plugin.Version)
	case "tag":
		referenceName = plumbing.NewTagReferenceName(plugin.Version)
	default:
		return fmt.Errorf("unknown version type for plugin '%+v': %s", plugin, plugin.VersionType)
	}

	cloneOptions := git.CloneOptions{
		URL:           plugin.Repository,
		ReferenceName: referenceName,
	}
	pullOptions := git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: referenceName,
		SingleBranch:  true,
	}

	repoPath := PluginsPath + plugin.ID

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Info("Repository does not exist, cloning it", "repository", plugin.Repository)
		err = CloneRepository(plugin, cloneOptions)
		if err != nil {
			return fmt.Errorf("error while cloning %v: %w", plugin.ID, err)
		}
	} else {
		checkoutOptions := git.CheckoutOptions{
			Branch: referenceName,
		}

		if err := Checkout(plugin, checkoutOptions); err != nil {
			return fmt.Errorf("error checking out branch %v: %w", referenceName, err)
		}

		log.Info("Repository exists, pulling it",
			slog.Group("plugin", "id", plugin.ID, "name", plugin.Name),
		)
		if err := PullRepository(plugin, pullOptions); err != nil {
			return fmt.Errorf("error pulling plugin '%+v': %w", plugin, err)
		}
	}

	return nil
}
