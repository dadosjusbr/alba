package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "alba"
const collector = "collector"

//Collector represents the information needed for frequent data collection operation
type Collector struct {
	ID                 string    `bson:"id"`                   // Initials entity like 'trt13'.
	Entity             string    `bson:"entity"`               // Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'.
	City               string    `bson:"city"`                 // City of the entity from which the collector extracts data.
	FU                 string    `bson:"fu"`                   // Federation unit of the entity from which the collector extracts data.
	UpdateDate         time.Time `bson:"update_date"`          // Last time the collector register has been updated.
	Path               string    `bson:"path"`                 // Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'.
	Frequency          int       `bson:"frequency"`            // Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'.
	StartDay           int       `bson:"start_day"`            // Day of the month for the collector execution. Values must be between 1 and 30.
	LimitMonthBackward int       `bson:"limit_month_backward"` // The limit month to which the collector must be executed in its historical execution.
	LimitYearBackward  int       `bson:"limit_year_backward"`  // The limit year until which the collector must be executed in its historical execution.
}

// InsertCollector insert an collector array
func InsertCollector(newCollector Collector) error {
	client, err := conect()
	if err != nil {
		return fmt.Errorf("connect error: %q", err)
	}

	database := client.Database(database)
	collectorCollection := database.Collection(collector)
	_, err = collectorCollection.InsertOne(context.TODO(), newCollector)
	if err != nil {
		return fmt.Errorf("insert error: %q", err)
	}

	err = disconect(client)
	if err != nil {
		return fmt.Errorf("disconect error: %q", err)
	}

	return nil
}

func conect() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB")
	if uri == "" {
		return nil, fmt.Errorf("error trying get environment variable:%q", errors.New("$MONGODB is empty"))
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error trying to connect:%q", err)
	}

	return client, nil
}

func disconect(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}

	return nil
}
