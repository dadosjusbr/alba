package mcollector

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "alba"
const collector = "collector"
const uri = "mongodb://root:example@mongo:27017"

//Collector represents the information needed for frequent data collection operation
type Collector struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Name               string             `bson:"name"`
	Entity             string             `bson:"entity"`
	City               string             `bson:"city"`
	Uf                 string             `bson:"uf"`
	UpdateDate         time.Time          `bson:"update_date"`
	Path               string             `bson:"path"`
	IDVersion          string             `bson:"id_version"`
	Frequency          int                `bson:"frequency"`
	StartDay           int                `bson:"start_day"`
	LimitMonthBackward int                `bson:"limit_month_backward"`
	LimitYearBackward  int                `bson:"limit_year_backward"`
}

// InsertCollector insert an collector array
func InsertCollector(collectors []interface{}) error {

	client, err := conect()
	if err != nil {
		return fmt.Errorf("connect error: %q", err)
	}

	database := client.Database(database)
	collectorCollection := database.Collection(collector)
	res, err := collectorCollection.InsertMany(context.TODO(), collectors)
	if err != nil {
		return fmt.Errorf("insert error: %q", err)
	}
	fmt.Println("inserted an array of documents: ", res.InsertedIDs)

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
