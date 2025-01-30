package pluginmanager

import (
	"log"
	"os"

	"github.com/epos-eu/converter-routine/dao/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CloneOrPull clones or pulls the software source code repositories
// the branch parameter determines whether to consider the software version as a branch or a tag
func CloneOrPull(sscs []model.Softwaresourcecode, branch bool) []model.Softwaresourcecode {
	installedRepos := make([]model.Softwaresourcecode, 0, len(sscs))

	// Iterate over each software source code object
	for _, obj := range sscs {
		// Determine the reference name based on the provided options
		var referenceName plumbing.ReferenceName
		if branch {
			referenceName = plumbing.NewBranchReferenceName(obj.Softwareversion)
		} else {
			referenceName = plumbing.NewTagReferenceName(obj.Softwareversion)
		}

		// Define clone and pull options
		cloneOptions := git.CloneOptions{
			URL:           obj.Coderepository,
			ReferenceName: referenceName,
		}
		pullOptions := git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: referenceName,
			SingleBranch:  true,
		}

		// Construct the repository path using the instance ID
		repoPath := PluginsPath + obj.InstanceID

		// Check if the repository directory exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			log.Printf("Repository %v does not exist, cloning...", obj.UID)
			// If the repository does not exist, clone it
			err = CloneRepository(obj, cloneOptions)
			// If there was an error cloning
			if err != nil {
				log.Printf("Error while cloning %v: %v", obj.UID, err)
				// don't add this to the installed repositories
				continue
			}
		} else {
			// Define checkout options
			checkoutOptions := git.CheckoutOptions{
				Branch: referenceName,
			}

			// Checkout the specified branch
			if err := Checkout(obj, checkoutOptions); err != nil {
				log.Printf("Error checking out branch %v: %v\n", referenceName, err)
			}

			log.Printf("Repository %v exists, pulling...\n", obj.UID)
			// If the repository exists, attempt to pull the latest changes
			if err := PullRepository(obj, pullOptions); err != nil {
				log.Printf("Error pulling: %v\n", err)
			}
		}

		// If we get here, it means it was successfully cloned/pulled
		installedRepos = append(installedRepos, obj)
	}

	return installedRepos
}
