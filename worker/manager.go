package worker

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

// CloneRepository is responsible for get the latest code version of pipeline repository.
// From the environment variable BASEDIR creates and returns the DefaultBaseDir for the pipeline.
func CloneRepository(repo string) (string, error) {
	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		return "", fmt.Errorf("error cloning the repository. BASEDIR env var can not be empty")
	}

	defaultBaseDir := fmt.Sprintf("%s/%s", baseDir, repo)
	if err := os.RemoveAll(defaultBaseDir); err != nil {
		return "", fmt.Errorf("error cloning the repository. error removing previous directory: %q", err)
	}

	url := fmt.Sprintf("https://%s", repo)
	log.Printf("Cloning the repository [%s] into [%s]\n", url, defaultBaseDir)
	_, err := git.PlainClone(defaultBaseDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", fmt.Errorf("error cloning the repository: %q", err)
	}

	return defaultBaseDir, nil
}
