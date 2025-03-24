package pluginmanager

import (
	"errors"

	"github.com/epos-eu/converter-routine/dao/model"
	"gopkg.in/src-d/go-git.v4"
)

func PullRepository(plugin model.Plugin, options git.PullOptions) error {
	// Open the given repository
	r, err := git.PlainOpen(PluginsPath + plugin.ID)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes
	err = w.Pull(&options)
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			// log.Println("Already up to date")
		} else {
			return err
		}
	}
	return nil
}
