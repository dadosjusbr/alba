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

func conect() *mongo.Client{
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

func insertCollector(client *mongo.Client){
	collection := client.Database(database).Collection(collection)
	res, err := collection.InsertOne(context.TODO(), bson.M{"name": "trt13", "path": "coletores/trt13/trt13"})
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", res.InsertedID)	
}

func disconect(client *mongo.Client){
	err := client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	fmt.Println("Connection to MongoDB closed.")
}