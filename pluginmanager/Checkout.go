package pluginmanager

import (
	"github.com/epos-eu/converter-routine/orms"
	"gopkg.in/src-d/go-git.v4"
)

func Checkout(obj orms.SoftwareSourceCode, options git.CheckoutOptions) error {
	// Open the given repository
	r, err := git.PlainOpen(PluginsPath + obj.GetInstanceID())
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
