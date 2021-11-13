package mongodb

import (
	"context"
	"example/hello/project/internal/store"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client

	songs   store.SongsRepository
	artists store.ArtistsRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(uri string) error {
	// set client options
	clientOptions := options.Client().ApplyURI(uri)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	log.Println("connected to MongoDB!")

	db.client = client
	return nil
}

func (db *DB) Close() error {
	// disconnecting
	err := db.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	log.Println("connection to MongoDB closed")
	return nil
}
