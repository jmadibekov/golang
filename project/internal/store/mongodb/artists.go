package mongodb

import (
	"context"
	"errors"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (db DB) Artists() store.ArtistsRepository {
	if db.artists == nil {
		artists, err := NewArtistsRepository(db.client)

		if err != nil {
			log.Fatalf("got an error while creating a %v collection with constraints. Err: %v", "artists", err)
			return nil
		}
		db.artists = artists
	}

	return db.artists
}

type ArtistsRepository struct {
	client *mongo.Client

	collection *mongo.Collection
}

func NewArtistsRepository(client *mongo.Client) (store.ArtistsRepository, error) {
	// if either the database or the collection doesn't exist, the following line will create them
	artistsCollection := client.Database("lostify").Collection("artists")
	// creating an index so that `ID` field is unique
	// read more: https://stackoverflow.com/questions/55921098/making-a-unique-field-in-mongo-go-driver
	_, err := artistsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return &ArtistsRepository{
		client:     client,
		collection: client.Database("lostify").Collection("artists"),
	}, nil
}

func (c ArtistsRepository) Create(ctx context.Context, artist *models.Artist) error {
	insertResult, err := c.collection.InsertOne(ctx, artist)
	if err != nil {
		return err
	}

	log.Println("inserted an artist:", insertResult.InsertedID)

	return nil
}

func (c ArtistsRepository) All(ctx context.Context) ([]*models.Artist, error) {
	// pass these options to the Find method
	findOptions := options.Find()
	// page size is set to 10
	findOptions.SetLimit(10)

	// here's an array in which you can store the decoded documents
	var artists []*models.Artist

	// passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := c.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	// finding multiple documents returns a cursor
	// iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var artist models.Artist
		err := cur.Decode(&artist)
		if err != nil {
			return nil, err
		}

		artists = append(artists, &artist)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// close the cursor once finished
	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("found multiple artists (array of pointers): %+v\n", artists)

	return artists, nil
}

func (c ArtistsRepository) ByID(ctx context.Context, id int) (*models.Artist, error) {
	// create a value into which the result can be decoded
	var artist models.Artist

	filter := bson.M{"id": id}

	if err := c.collection.FindOne(ctx, filter).Decode(&artist); err != nil {
		return nil, err
	}

	log.Printf("found an artist: %+v\n", artist)

	return &artist, nil
}

func (c ArtistsRepository) Update(ctx context.Context, artist *models.Artist) error {
	filter := bson.M{"id": artist.ID}

	updateResult, err := c.collection.ReplaceOne(ctx, filter, artist)
	if err != nil {
		return err
	}

	log.Printf("matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	if updateResult.MatchedCount != 1 {
		return errors.New("either no or more than one artist has been matched")
	}

	return nil
}

func (c ArtistsRepository) Delete(ctx context.Context, id int) error {
	filter := bson.M{"id": id}

	deleteResult, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Printf("deleted %v documents in the artists collection\n", deleteResult.DeletedCount)

	if deleteResult.DeletedCount != 1 {
		return errors.New("either no or more than one artist has been matched")
	}

	return nil
}
