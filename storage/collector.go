package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "alba"
const collector = "collector"
const uri = "mongodb://root:example@mongo:27017"

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
	res, err := collectorCollection.InsertOne(context.TODO(), newCollector)
	if err != nil {
		return fmt.Errorf("insert error: %q", err)
	}
	fmt.Println("inserted an array of documents: ", res.InsertedID)

	disconect := disconect(client)
	if disconect != nil {
		return fmt.Errorf("disconect error: %q", disconect)
	}

	return nil
}

func conect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return client, err
	}
	fmt.Println("connected to MongoDB!")

	return client, nil
}

func disconect(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}
	fmt.Println("Connection to MongoDB closed.")

	return nil
}
