package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/alba/git"
	"github.com/dadosjusbr/alba/storage"
)

type dbInterface interface {
	GetPipelinesByDay(day int) ([]storage.Pipeline, error)
	InsertExecution(e storage.Execution) error
	GetLastExecutionsForAllPipelines() ([]storage.Execution, error)
	GetLastExecutionsByPipelineID(limit, id int) ([]storage.Execution, error)
}

func main() {
	var pipelines []storage.Pipeline
	var finalPipelines []storage.Pipeline

	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}

	errorLimitStr := os.Getenv("ERROR_LIMIT")
	if uri == "" {
		log.Fatal("error trying get environment variable: $ERROR_LIMIT is empty")
	}

	errorLimit, err := strconv.Atoi(errorLimitStr)
	if err != nil {
		log.Fatalf("error trying convert variable $ERROR_LIMIT: %q", err)
	}

	dbClient, err := storage.NewDBClient(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = dbClient.Connect(); err != nil {
		log.Fatal(err)

	}

	// Setup for send email
	emailSender := os.Getenv("EMAIL_SENDER")
	if emailSender == "" {
		log.Fatal("setup error sending email. EMAIL_SENDER env var can not be empty")
	}

	emailPassword := os.Getenv("EMAIL_SENDER_PASSWORD")
	if emailPassword == "" {
		log.Fatal("setup error sending email. EMAIL_SENDER_PASSWORD env var can not be empty")
	}

	emailReceiver := os.Getenv("EMAIL_RECEIVERS")

	pipelines, err = getPipelinesToExecuteToday(dbClient)
	if err != nil {
		log.Fatal(err)
	}
	finalPipelines = append(finalPipelines, pipelines...)

	// pipelines, err = getPipelinesThatFailed(dbClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// TODO: merge pipelines a fim de não repetir os ids
	// final_pipelines = append(final_pipelines, pipelines...)

	// pipelines, err = getPipelinesForCompleteHistory(dbClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	toExecuteNow := prioritizeAndLimit(dbClient, finalPipelines, errorLimit)

	for _, p := range toExecuteNow {
		err := run(emailSender, emailPassword, emailReceiver, p, dbClient)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dbClient.Disconnect()
}

func run(emailSender, emailPassword, emailReceiver string, p storage.Pipeline, db *storage.DBClient) error {
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

	if emailReceiver != "" {
		if err := sendEmail(emailSender, emailPassword, emailReceiver, e.Entity, e.PipelineResult.Status); err != nil {
			return fmt.Errorf("error after running pipeline. %q", err)

		}
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

func prioritizeAndLimit(db dbInterface, list []storage.Pipeline, limit int) []storage.Pipeline {

	return list
}

func getPipelinesToExecuteToday(db dbInterface) ([]storage.Pipeline, error) {
	results, err := db.GetPipelinesByDay(time.Now().Day())
	if err != nil {
		return nil, fmt.Errorf("error getting pipelines by day: %q", err)
	}

	return results, nil
}

// Apenas as execuções que devem acontecer por causa do mecanismo de
// tolerância à falhas.
func getPipelinesThatFailed(db dbInterface) []storage.Pipeline {

	return nil
}

// Apenas execuções de devem acontecer para completar o histórico. Devemos
// ignorar casos em que já houve tentativa de execução, quer seja sucesso ou falha.
func getPipelinesForCompleteHistory(db dbInterface) []storage.Pipeline {

	return nil
}

func sendEmail(sender, password, receiver, entity, status string) error {
	// TODO: migrar para templates
	// TODO: modificar email para apontar para URL que a pessoa operadora possa ver o erro que aconteceu (a URL base do alba também deve ser uma variável de ambiente)

	auth := smtp.PlainAuth("", sender, password, "smtp.gmail.com")
	receivers := strings.Split(receiver, ",")

	message := []byte(fmt.Sprintf("To: %v \r\n"+
		"Subject: [DadosJusBR: Alba] Extraímos novos dados! \r\n"+
		"\r\n"+
		"Olá, sou a Alba e acabei de executar o pipeline para o órgão: %s com status: %s!\n"+
		"Acompanhe mais sobre o meu trabalho no site: https://dadosjusbr.org/", receivers, entity, status))

	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, receivers, message)
	if err != nil {
		return fmt.Errorf("error sending email: %q", err)
	}

	return nil
}
