package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const database = "alba"
const collector = "collector"

//Collector represents the information needed for frequent data collection operation
type Collector struct {
	ID                 string    `bson:"id, omitempty" json:"id"`                                     // Initials entity like 'trt13'.
	Entity             string    `bson:"entity, omitempty" json:"entity"`                             // Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'.
	City               string    `bson:"city, omitempty" json:"city"`                                 // City of the entity from which the collector extracts data.
	FU                 string    `bson:"fu, omitempty" json:"fu"`                                     // Federation unit of the entity from which the collector extracts data.
	UpdateDate         time.Time `bson:"update-date, omitempty" json:"update-date"`                   // Last time the collector register has been updated.
	Path               string    `bson:"path, omitempty" json:"path"`                                 // Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'.
	Frequency          int       `bson:"frequency, omitempty" json:"frequency"`                       // Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'.
	StartDay           int       `bson:"start-day, omitempty" json:"start-day"`                       // Day of the month for the collector execution. Values must be between 1 and 30.
	LimitMonthBackward int       `bson:"limit-month-backward, omitempty" json:"limit-month-backward"` // The limit month to which the collector must be executed in its historical execution.
	LimitYearBackward  int       `bson:"limit-year-backward, omitempty" json:"limit-year-backward"`   // The limit year until which the collector must be executed in its historical execution.
}

// InsertCollector insert an collector array
func InsertCollector(newCollector Collector) error {
	client, err := conect()
	if err != nil {
		return fmt.Errorf("connect error: %q", err)
	}

	collectorCollection := client.Database(database).Collection(collector)
	setIndexesCollector(collectorCollection)
	_, err = collectorCollection.InsertOne(context.TODO(), newCollector)
	if err != nil {
		return fmt.Errorf("insert error: %q", err)
	}

	disconect := disconect(client)
	if disconect != nil {
		return fmt.Errorf("disconect error: %q", disconect)
	}

	return nil
}

func setIndexesCollector(collectorCollection *mongo.Collection) error {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	indexes := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "path", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	}

	namesIndexes, err := collectorCollection.Indexes().CreateMany(context.Background(), indexes, opts)
	if err != nil {
		return fmt.Errorf("create index error: %q", err)
	}
	fmt.Println(namesIndexes)

	return nil
}

// GetCollectors return all collectors in the database
func GetCollectors() ([]*Collector, error) {
	var collectors []*Collector

	client, err := conect()
	if err != nil {
		return collectors, fmt.Errorf("connect error: %q", err)
	}

	collectorCollection := client.Database(database).Collection(collector)
	itens, err := collectorCollection.Find(context.TODO(), bson.D{{}})
	for itens.Next(context.TODO()) {
		var item Collector
		err := itens.Decode(&item)
		if err != nil {
			return collectors, fmt.Errorf("decode error in collector: %q", err)
		}

		collectors = append(collectors, &item)
	}
	itens.Close(context.TODO())

	disconect := disconect(client)
	if disconect != nil {
		return collectors, fmt.Errorf("disconect error: %q", disconect)
	}

	return collectors, nil
}

// GetCollector find and return a collector from the path
func GetCollector(path string) (Collector, error) {

	return Collector{}, nil
}

func conect() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return client, err
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
