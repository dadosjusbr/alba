package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/dadosjusbr/executor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const database = "alba"
const pipelineCollection = "pipeline"
const executionCollection = "execution"

// Pipeline represents the information needed for frequent data collection operation.
type Pipeline struct {
	Pipeline           executor.Pipeline `bson:"pipeline, omitempty" json:"pipeline"`                         // Represents the sequence of stages for data release.
	ID                 string            `bson:"id, omitempty" json:"id"`                                     // Initials entity like 'trt13'.
	Entity             string            `bson:"entity, omitempty" json:"entity"`                             // Entity from which the pipeline extracts data like 'Tribunal Regional do Trabalho 13° Região'.
	City               string            `bson:"city, omitempty" json:"city"`                                 // City of the entity from which the pipeline extracts data.
	FU                 string            `bson:"fu, omitempty" json:"fu"`                                     // Federation unit of the entity from which the pipeline extracts data.
	Repo               string            `bson:"repo, omitempty" json:"repo"`                                 // Central pipeline repository. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'.
	Frequency          int               `bson:"frequency, omitempty" json:"frequency"`                       // Frequency of the pipeline execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'.
	StartDay           int               `bson:"start-day, omitempty" json:"start-day"`                       // Day of the month for the pipeline execution. Values must be between 1 and 30.
	LimitMonthBackward int               `bson:"limit-month-backward, omitempty" json:"limit-month-backward"` // The limit month to which the pipeline must be executed in its historical execution.
	LimitYearBackward  int               `bson:"limit-year-backward, omitempty" json:"limit-year-backward"`   // The limit year until which the pipeline must be executed in its historical execution.
	UpdateDate         time.Time         `bson:"update-date, omitempty" json:"update-date"`                   // Last time the pipeline register has been updated.
}

// DBClient represents a mongodb client instance.
type DBClient struct {
	mgoClient *mongo.Client
}

// NewDBClient return a DBCLient.
func NewDBClient(uri string) (*DBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("new dbclient error. error creating new client: %q", err)
	}
	return &DBClient{mgoClient: client}, nil
}

// Connect makes the connection and setup of database.
func (c *DBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("connect error. error trying to connect: %q", err)
	}

	// Check if alba database exist.
	albaExists, err := c.mgoClient.ListDatabaseNames(ctx, bson.D{{Key: "name", Value: database}})
	if err != nil {
		return fmt.Errorf("connect error. error when listing database names: %q", err)
	}

	if len(albaExists) == 0 { // Database setup.
		collection := c.mgoClient.Database(database).Collection(pipelineCollection)
		if err := setIndexesPipeline(collection); err != nil {
			return fmt.Errorf("connect error. set indexes error in collection %s: %q", pipelineCollection, err)
		}
	}

	return nil
}

func setIndexesPipeline(pipeline *mongo.Collection) error {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	indexes := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "repo", Value: bsonx.Int32(1)}, {Key: "id", Value: bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := pipeline.Indexes().CreateOne(context.Background(), indexes, opts); err != nil {
		return fmt.Errorf("create index error: %q", err)
	}

	return nil
}

// InsertPipeline insert a pipeline in the database.
func (c *DBClient) InsertPipeline(newPipeline Pipeline) error {
	collection := c.mgoClient.Database(database).Collection(pipelineCollection)
	if _, err := collection.InsertOne(context.TODO(), newPipeline); err != nil {
		return fmt.Errorf("insert error: %q", err)
	}

	return nil
}

// GetPipelines return all pipelines in the database.
func (c *DBClient) GetPipelines() ([]Pipeline, error) {
	var pipelines []Pipeline

	collection := c.mgoClient.Database(database).Collection(pipelineCollection)
	itens, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []Pipeline{}, nil
		}
		return nil, fmt.Errorf("error getting pipelines. Find error: %q", err)
	}

	for itens.Next(context.Background()) {
		var item Pipeline
		if err := itens.Decode(&item); err != nil {
			return nil, fmt.Errorf("error getting pipelines. Decode error: %q", err)
		}
		pipelines = append(pipelines, item)
	}
	itens.Close(context.Background())

	return pipelines, nil
}

// Disconnect makes the database disconnection.
func (c *DBClient) Disconnect() error {
	if err := c.mgoClient.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}

	return nil
}
