package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dadosjusbr/alba"
	"github.com/dadosjusbr/alba/storage"
)

func main() {
	var pipelines []storage.Pipeline

	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}

	dbClient, err := storage.NewDBClient(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = dbClient.Connect(); err != nil {
		log.Fatal(err)

	}

	pipelines = append(pipelines, getPipelinesToExecuteNow(dbClient)...)
	pipelines = append(pipelines, getPipelinesThatFailed(dbClient)...)
	pipelines = append(pipelines, getPipelinesForCompleteHistory(dbClient)...)

	// Algoritmo: shuffle na lista + cap
	toExecuteNow := prioritizeAndLimit(pipelines)
	for _, p := range toExecuteNow {
		err := run(p, dbClient)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dbClient.Disconnect()
}

func run(p storage.Pipeline, db *storage.DBClient) error {
	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		return fmt.Errorf("error running pipeline. BASEDIR env var can not be empty")
	}

	month := os.Getenv("MONTH")
	if month == "" {
		return fmt.Errorf("error running pipeline. MONTH env var can not be empty")
	}

	year := os.Getenv("YEAR")
	if year == "" {
		return fmt.Errorf("error running pipeline. YEAR env var can not be empty")
	}

	var commit string
	p.Pipeline.DefaultBaseDir = fmt.Sprintf("%s/%s", baseDir, p.Repo)
	commit, err := alba.CloneRepository(p.Pipeline.DefaultBaseDir, fmt.Sprintf("https://%s", p.Repo))
	if err != nil {
		return fmt.Errorf("error running pipeline: %q", err)
	}

	defaultBuildEnv := map[string]string{
		"GIT_COMMIT": commit,
	}

	defaultRunEnv := map[string]string{
		"OUTPUT_FOLDER": "/output",
		"MONTH":         month,
		"YEAR":          year,
	}

	for pos, stage := range p.Pipeline.Stages {
		stage.BuildEnv = mergeEnv(defaultBuildEnv, stage.BuildEnv)
		stage.RunEnv = mergeEnv(defaultRunEnv, stage.RunEnv)
		p.Pipeline.Stages[pos] = stage
	}

	result, _ := p.Pipeline.Run()
	e := storage.Execution{
		PipelineResult: result,
		Entity:         p.Entity,
		ID:             p.ID,
	}
	db.InsertExecution(e)
	return nil
}

func mergeEnv(defaultEnv, stageEnv map[string]string) map[string]string {
	env := make(map[string]string)

	for k, v := range defaultEnv {
		env[k] = v
	}
	for k, v := range stageEnv {
		env[k] = v
	}
	return env
}

func prioritizeAndLimit(list []storage.Pipeline) []storage.Pipeline {

	return nil
}

// Assumindo que o passado não interessa, quais pipelines devem ser
// executados no dia/hora atual
func getPipelinesToExecuteNow(db *storage.DBClient) []storage.Pipeline {

	return nil
}

// Apenas as execuções que devem acontecer por causa do mecanismo de
// tolerância à falhas.
func getPipelinesThatFailed(db *storage.DBClient) []storage.Pipeline {

	return nil
}

// Apenas execuções de devem acontecer para completar o histórico. Devemos
// ignorar casos em que já houve tentativa de execução, quer seja sucesso ou falha.
func getPipelinesForCompleteHistory(db *storage.DBClient) []storage.Pipeline {

	return nil
}
