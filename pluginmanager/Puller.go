package pluginmanager

import (
	"errors"
	"github.com/epos-eu/converter-routine/orms"
	"gopkg.in/src-d/go-git.v4"
)

func PullRepository(obj orms.SoftwareSourceCode, options git.PullOptions) error {
	// Open the given repository
	r, err := git.PlainOpen(PluginsPath + obj.GetInstance_id())
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
