// following great this tutorial: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// set client options
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://dbUser:dbUserPassword@cluster0.99ocj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type of client is %T\n", client)

	// check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// accessing 'trainers' collection in MongoDB 'test' database
	collection := client.Database("test").Collection("trainers")
	fmt.Printf("Type of collection is %T\n", collection)

	// ash := Trainer{"Ash", 10, "Pallet Town"}
	// misty := Trainer{"Misty", 10, "Cerulean City"}
	// brock := Trainer{"Brock", 15, "Pewter City"}

	// // inserting a single document, in this case ash
	// insertResult, err := collection.InsertOne(ctx, ash)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// // inserting multiple documents
	// trainers := []interface{}{misty, brock}

	// insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	// disconnecting
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
