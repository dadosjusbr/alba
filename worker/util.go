package worker

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/dadosjusbr/alba/storage"
)

// CloneRepository is responsible for get the latest code version of pipeline repository.
// Creates and returns the DefaultBaseDir for the pipeline and the latest commit in the repository.
func CloneRepository(repo string) (string, string, error) {
	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		return "", "", fmt.Errorf("error cloning the repository. BASEDIR env var can not be empty")
	}

	defaultBaseDir := fmt.Sprintf("%s/%s", baseDir, repo)
	if err := os.RemoveAll(defaultBaseDir); err != nil {
		return "", "", fmt.Errorf("error cloning the repository. error removing previous directory: %q", err)
	}

	url := fmt.Sprintf("https://%s", repo)
	log.Printf("Cloning the repository [%s] into [%s]\n", url, defaultBaseDir)
	r, err := git.PlainClone(defaultBaseDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", "", fmt.Errorf("error cloning the repository: %q", err)
	}

	ref, err := r.Head()
	if err != nil {
		return "", "", fmt.Errorf("error cloning the repository. error getting the HEAD reference of the repository: %q", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", "", fmt.Errorf("error cloning the repository. error getting the lattest commit of the repository: %q", err)
	}
	return defaultBaseDir, fmt.Sprintf("%s", commit.Hash), nil
}

// SetupDadosjusBR TODO: Função que reúne as regras de negócio para um pipeline DadosJusBR,
// como configuração das variáveis commit, mes e ano
func SetupDadosjusBR(p storage.Pipeline, month, year string) (storage.Pipeline, error) {
	return storage.Pipeline{}, nil
}
