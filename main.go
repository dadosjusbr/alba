package main

import (
	"fmt"
	"log"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	insert_collector()
}

func insert_collector(){

	client := conect()
	collection := client.Database("alba").Collection("collector")
	res, err := collection.InsertOne(context.TODO(), bson.M{"name": "trt13", "path": "coletores/trt13/trt13"})
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Inserted a single document: ", res.InsertedID)	

	disconect(client)
}

func conect() *mongo.Client{

	clientOptions := options.Client().ApplyURI("mongodb://root:example@mongo:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func disconect(client *mongo.Client){

	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}