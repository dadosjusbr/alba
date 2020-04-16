package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const database = "alba"
const collectorCollection = "collector"

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

// InsertCollector insert a collector
func InsertCollector(newCollector Collector) error {
	client, err := connect()
	if err != nil {
		return fmt.Errorf("connect error: %q", err)
	}

	collectorC := client.Database(database).Collection(collectorCollection)
	if collectorC == nil {
		return fmt.Errorf("error in retrive collection")
	}

	if _, err = collectorC.InsertOne(context.TODO(), newCollector); err != nil {
		return fmt.Errorf("insert error: %q", err)
	}

	if err = disconnect(client); err != nil {
		return fmt.Errorf("disconnect error: %q", err)
	}

	return nil
}

// GetCollectors return all collectors in the database
func GetCollectors() ([]Collector, error) {
	var collectors []Collector

	client, err := connect()
	if err != nil {
		return nil, fmt.Errorf("connect error: %q", err)
	}

	collectorC := client.Database(database).Collection(collectorCollection)
	itens, err := collectorC.Find(context.TODO(), bson.D{{}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("Find error: %q", err)
	}

	for itens.Next(context.Background()) {
		var item Collector
		err := itens.Decode(&item)
		if err != nil {
			return nil, fmt.Errorf("decode error in collector: %q", err)
		}
		collectors = append(collectors, item)
	}
	itens.Close(context.Background())

	err = disconnect(client)
	if err != nil {
		return nil, fmt.Errorf("disconnect error: %q", err)
	}

	collectorsJSON, err := json.Marshal(collectors)
	if err != nil {
		return nil, fmt.Errorf("json encodind error: %q", err)
	}

	return collectorsJSON, nil
}

func connect() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB")
	if uri == "" {
		return nil, fmt.Errorf("error trying get environment variable:%q", errors.New("$MONGODB is empty"))
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error trying to connect:%q", err)
	}
	//Check if alba database exist
	results, err := client.ListDatabaseNames(context.TODO(), bson.D{{Key: "name", Value: database}})
	if err != nil {
		return nil, fmt.Errorf("error when listing database names: %q", err)
	}

	if len(results) == 0 { //First executition for alba database and setup
		err := setupDB(client)
		if err != nil {
			return nil, fmt.Errorf("setup database error: %q", err)
		}
	}

	return client, nil
}

//setupDB creates the collections and indexes
func setupDB(client *mongo.Client) error {
	collectorC := client.Database(database).Collection(collectorCollection)
	if collectorC == nil {
		return fmt.Errorf("error in create collection: %q", collectorCollection)
	}

	err := setIndexesCollector(collectorC)
	if err != nil {
		return fmt.Errorf("set indexes error in collection: %q", collectorCollection)
	}

	return nil
}

func setIndexesCollector(collectorC *mongo.Collection) error {
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

	_, err := collectorC.Indexes().CreateMany(context.Background(), indexes, opts)
	if err != nil {
		return fmt.Errorf("create index error: %q", err)
	}

	return nil
}

func disconnect(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}

	return nil
}
