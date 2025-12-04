package pluginmanager

import (
	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/go-git/go-git/v5"
)

func Checkout(plugin model.Plugin, options git.CheckoutOptions) error {
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

	// Checkout the branch
	err = w.Checkout(&options)
	if err != nil {
		return err
	}

	return nil
}
