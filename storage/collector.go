package storage

import (
	"context"
	"fmt"
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

//DBClient represents a mongodb client instance
type DBClient struct {
	mgoClient *mongo.Client
}

//NewDBClient return a DBCLient
func NewDBClient(uri string) (*DBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("new dbclient error. error creating new client: %q", err)
	}
	return &DBClient{mgoClient: client}, nil
}

//Connect makes the connection to the database
func (c *DBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("connect error. error trying to connect:%q", err)
	}

	//Check if alba database exist
	results, err := c.mgoClient.ListDatabaseNames(ctx, bson.D{{Key: "name", Value: database}})
	if err != nil {
		return fmt.Errorf("connect error. error when listing database names: %q", err)
	}

	if len(results) == 0 { //First execution for alba database and setup
		collection := c.mgoClient.Database(database).Collection(collectorCollection)
		if err := setIndexesCollector(collection); err != nil {
			return fmt.Errorf("connect error. set indexes error in collection: %q", collectorCollection)
		}
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

	if _, err := collectorC.Indexes().CreateMany(context.Background(), indexes, opts); err != nil {
		return fmt.Errorf("create index error: %q", err)
	}

	return nil
}

//InsertCollector insert a collector
func (c *DBClient) InsertCollector(newCollector Collector) error {
	collection := c.mgoClient.Database(database).Collection(collectorCollection)
	if _, err := collection.InsertOne(context.TODO(), newCollector); err != nil {
		return fmt.Errorf("insert error: %q", err)
	}

	return nil
}

// GetCollectors return all collectors in the database
func (c *DBClient) GetCollectors() ([]Collector, error) {
	var collectors []Collector

	collection := c.mgoClient.Database(database).Collection(collectorCollection)
	itens, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []Collector{}, nil
		}
		return nil, fmt.Errorf("error getting collectors. Find error: %q", err)
	}

	for itens.Next(context.Background()) {
		var item Collector
		if err := itens.Decode(&item); err != nil {
			return nil, fmt.Errorf("error getting collectors. Decode error: %q", err)
		}
		collectors = append(collectors, item)
	}
	itens.Close(context.Background())

	return collectors, nil
}

//Disconnect makes the disconnection to the database
func (c *DBClient) Disconnect() error {
	if err := c.mgoClient.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}

	return nil
}
