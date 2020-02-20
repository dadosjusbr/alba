package main

import (
	"fmt"
	"log"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "alba"
const collection = "collector"
const uri = "mongodb://root:example@mongo:27017"

func main() {

	client, err := conect()
	if err != nil {
		log.Fatal(err)
	}
	if insertCollector(client) != nil {
		log.Fatal("Insert error")
	}
	if disconect(client) != nil {
		log.Fatal("Disconect error")	
	}
}

func conect() (*mongo.Client, error){
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return client, err
	}
	fmt.Println("Connected to MongoDB!")

	return client, nil
}

func insertCollector(client *mongo.Client) error{
	collection := client.Database(database).Collection(collection)
	res, err := collection.InsertOne(context.TODO(), bson.M{"name": "trt13", "path": "coletores/trt13/trt13"})
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)	

	return nil
}

func disconect(client *mongo.Client) error{
	err := client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error trying to disconnect:%q", err)
	}
	fmt.Println("Connection to MongoDB closed.")

	return nil
}