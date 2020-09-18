package alba

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

// CloneRepository is responsible for get the latest code version of pipeline repository.
// Creates and returns the DefaultBaseDir for the pipeline and the latest commit in the repository.
func CloneRepository(defaultBaseDir, url string) (string, error) {
	if err := os.RemoveAll(defaultBaseDir); err != nil {
		return "", fmt.Errorf("error cloning the repository. error removing previous directory: %q", err)
	}

	log.Printf("Cloning the repository [%s] into [%s]\n\n", url, defaultBaseDir)
	r, err := git.PlainClone(defaultBaseDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", fmt.Errorf("error cloning the repository: %q", err)
	}

	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("error cloning the repository. error getting the HEAD reference of the repository: %q", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", fmt.Errorf("error cloning the repository. error getting the lattest commit of the repository: %q", err)
	}
	return fmt.Sprintf("%s", commit.Hash), nil
}
