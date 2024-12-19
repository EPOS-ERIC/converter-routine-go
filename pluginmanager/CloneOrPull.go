package pluginmanager

import (
	"github.com/epos-eu/converter-routine/orms"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"os"
)

// CloneOrPull clones or pulls the software source code repositories
// the branch parameter determines whether to consider the software version as a branch or a tag
func CloneOrPull(sscs []orms.SoftwareSourceCode, branch bool) []orms.SoftwareSourceCode {
	installedRepos := make([]orms.SoftwareSourceCode, len(sscs))
	copy(installedRepos, sscs)

	// Iterate over each software source code object
	for i, obj := range sscs {
		// Determine the reference name based on the provided options
		var referenceName plumbing.ReferenceName
		if branch {
			referenceName = plumbing.NewBranchReferenceName(obj.GetSoftwareVersion())
		} else {
			referenceName = plumbing.NewTagReferenceName(obj.GetSoftwareVersion())
		}

		// Define clone and pull options
		cloneOptions := git.CloneOptions{
			URL:           obj.GetCodeRepository(),
			ReferenceName: referenceName,
		}
		pullOptions := git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: referenceName,
			SingleBranch:  true,
		}

		// Construct the repository path using the instance ID
		repoPath := PluginsPath + obj.GetInstanceID()

		// Check if the repository directory exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			log.Printf("Repository %v does not exist, cloning...", obj.UID)
			// If the repository does not exist, clone it
			err = CloneRepository(obj, cloneOptions)
			// If there was an error cloning
			if err != nil {
				// Remove from the installed repos
				installedRepos = append(installedRepos[:i], installedRepos[i+1:]...)
				log.Printf("Error while cloning %v: %v", obj.UID, err)
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
	}

	return installedRepos
}
