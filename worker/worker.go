package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/dadosjusbr/alba/git"
	"github.com/dadosjusbr/alba/storage"
)

func main() {
	var pipelines []storage.Pipeline
	var finalPipelines []storage.Pipeline

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

	pipelines, err = getPipelinesToExecuteToday(dbClient)
	if err != nil {
		log.Fatal(err)
	}
	finalPipelines = append(finalPipelines, pipelines...)

	// pipelines, err = getPipelinesThatFailed(dbClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// final_pipelines = append(final_pipelines, pipelines...)

	// pipelines, err = getPipelinesForCompleteHistory(dbClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// final_pipelines = append(final_pipelines, pipelines...)

	// Algoritmo: shuffle na lista + cap
	toExecuteNow := prioritizeAndLimit(finalPipelines)

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

	month := int(time.Now().Month())
	year := time.Now().Year()
	if month == 1 {
		month = 12
		year = year - 1
	} else {
		month = month - 1
	}

	var commit string
	p.Pipeline.DefaultBaseDir = fmt.Sprintf("%s/%s", baseDir, p.Repo)
	commit, err := git.CloneRepository(p.Pipeline.DefaultBaseDir, fmt.Sprintf("https://%s", p.Repo))
	if err != nil {
		return fmt.Errorf("error running pipeline: %q", err)
	}

	defaultBuildEnv := map[string]string{
		"GIT_COMMIT": commit,
	}

	defaultRunEnv := map[string]string{
		"OUTPUT_FOLDER": "/output",
		"MONTH":         strconv.Itoa(month),
		"YEAR":          strconv.Itoa(year),
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

	receiver := "default@dadosjusbr.com"
	if err := sendEmail(receiver, e.Entity, e.PipelineResult.Status); err != nil {
		return fmt.Errorf("error after running pipeline. %q", err)

	}

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

	return list
}

func getPipelinesToExecuteToday(db *storage.DBClient) ([]storage.Pipeline, error) {

	day := time.Now().Day()
	fmt.Println(day)
	results, err := db.GetPipelinesByDay(day)
	if err != nil {
		return nil, fmt.Errorf("error getting pipelines by day: %q", err)
	}

	return results, nil
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

func sendEmail(receiver, entity, status string) error {
	sender := os.Getenv("SENDER")
	if sender == "" {
		return fmt.Errorf("setup error sending email. SENDER env var can not be empty")
	}

	password := os.Getenv("SENDER_PASSWORD")
	if password == "" {
		return fmt.Errorf("setup error sending email. SENDER_PASSWORD env var can not be empty")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	auth := smtp.PlainAuth("", sender, password, smtpHost)

	receivers := []string{
		receiver,
	}

	message := []byte(fmt.Sprintf("To: %v \r\n"+
		"Subject: [DadosJusBR: Alba] Extraímos novos dados! \r\n"+
		"\r\n"+
		"Olá, sou a Alba e acabei de executar o pipeline para o órgão: %s com status: %s!\n"+
		"Acompanhe mais sobre o meu trabalho no site: https://dadosjusbr.org/", receivers, entity, status))

	err := smtp.SendMail(addr, auth, sender, receivers, message)
	if err != nil {
		return fmt.Errorf("error sending email: %q", err)
	}

	return nil
}
