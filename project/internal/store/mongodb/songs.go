// mongodb-go driver tutorial: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// great tutorial on updating in mongodb: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents

package mongodb

import (
	"context"
	"errors"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) Songs() store.SongsRepository {
	if db.songs == nil {
		songs, err := NewSongsRepository(db.client)

		if err != nil {
			log.Fatalf("got an error while creating a %v collection with constraints. Err: %v", "songs", err)
			return nil
		}
		db.songs = songs
	}

	return db.songs
}

type SongsRepository struct {
	client *mongo.Client

	collection *mongo.Collection
}

func NewSongsRepository(client *mongo.Client) (store.SongsRepository, error) {
	// if either the database or the collection doesn't exist, the following line will create them
	songsCollection := client.Database("lostify").Collection("songs")
	// creating an index so that `ID` field is unique
	// read more: https://stackoverflow.com/questions/55921098/making-a-unique-field-in-mongo-go-driver
	_, err := songsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return &SongsRepository{
		client:     client,
		collection: client.Database("lostify").Collection("songs"),
	}, nil
}

func (c SongsRepository) GetArtist(ctx context.Context, song *models.Song) (*models.Artist, error) {
	// checks whether artist with `song.ArtistID` exists in the database
	var artist models.Artist

	filter := bson.M{"id": song.ArtistID}

	artistsCollection := c.client.Database("lostify").Collection("artists")

	if err := artistsCollection.FindOne(ctx, filter).Decode(&artist); err != nil {
		return nil, err
	}

	log.Printf("found an artist with artistID %v: %+v\n", song.ArtistID, artist)

	return &artist, nil
}

func (c SongsRepository) Create(ctx context.Context, song *models.Song) error {
	_, err := c.GetArtist(ctx, song)
	if err != nil {
		return errors.New(fmt.Sprintf("artist with id %v doesn't exist", song.ArtistID))
	}

	insertResult, err := c.collection.InsertOne(ctx, song)
	if err != nil {
		return err
	}

	log.Println("inserted a song:", insertResult.InsertedID)

	return nil
}

func (c SongsRepository) All(ctx context.Context, filter *models.Filter) ([]*models.Song, error) {
	// pass these options to the Find method
	findOptions := options.Find()
	// page size is set to 10
	findOptions.SetLimit(10)

	// here's an array in which you can store the decoded documents
	var songs []*models.Song

	// bson.D{{}} corresponds to the filter that matches all documents in the collection
	bsonFilter := bson.D{{}}
	if filter.Query != nil {
		// finds any songs with title that contains filter.Query (case insensitive search)
		// https://docs.mongodb.com/v4.4/tutorial/query-documents/
		// https://stackoverflow.com/questions/3305561/how-to-query-mongodb-with-like
		bsonFilter = bson.D{
			{"title", primitive.Regex{Pattern: *filter.Query, Options: "i"}},
		}
	}
	cur, err := c.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, err
	}

	// finding multiple documents returns a cursor
	// iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var song models.Song
		err := cur.Decode(&song)
		if err != nil {
			return nil, err
		}

		songs = append(songs, &song)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// close the cursor once finished
	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("found multiple songs (array of pointers): %+v\n", songs)

	return songs, nil
}

func (c SongsRepository) ByID(ctx context.Context, id int) (*models.Song, error) {
	// create a value into which the result can be decoded
	var song models.Song

	filter := bson.M{"id": id}

	if err := c.collection.FindOne(ctx, filter).Decode(&song); err != nil {
		return nil, err
	}

	log.Printf("found a song: %+v\n", song)

	return &song, nil
}

func (c SongsRepository) Update(ctx context.Context, song *models.Song) error {
	_, err := c.GetArtist(ctx, song)
	if err != nil {
		return errors.New(fmt.Sprintf("artist with %+v doesn't exist", song.ArtistID))
	}

	filter := bson.M{"id": song.ID}

	updateResult, err := c.collection.ReplaceOne(ctx, filter, song)
	if err != nil {
		return err
	}

	log.Printf("matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	if updateResult.MatchedCount != 1 {
		return errors.New("either no or more than one song has been matched")
	}

	return nil
}

func (c SongsRepository) Delete(ctx context.Context, id int) error {
	filter := bson.M{"id": id}

	deleteResult, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Printf("deleted %v documents in the songs collection\n", deleteResult.DeletedCount)

	if deleteResult.DeletedCount != 1 {
		return errors.New("either no or more than one song has been matched")
	}

	return nil
}
