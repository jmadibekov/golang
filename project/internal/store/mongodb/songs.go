// mongodb-go driver tutorial: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// great tutorial on updating in mongodb: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents

package mongodb

import (
	"context"
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) Songs() store.SongsRepository {
	if db.songs == nil {
		db.songs = NewSongsRepository(db.client)
	}

	return db.songs
}

type SongsRepository struct {
	client *mongo.Client

	collection *mongo.Collection
}

func NewSongsRepository(client *mongo.Client) store.SongsRepository {
	return &SongsRepository{
		client:     client,
		collection: client.Database("lostify").Collection("songs"),
	}
}

func (c SongsRepository) Create(ctx context.Context, song *models.Song) error {
	insertResult, err := c.collection.InsertOne(ctx, song)
	if err != nil {
		return err
	}

	log.Println("inserted a song:", insertResult.InsertedID)

	return nil
}

func (c SongsRepository) All(ctx context.Context) ([]*models.Song, error) {
	// pass these options to the Find method
	findOptions := options.Find()
	// page size is set to 10
	findOptions.SetLimit(10)

	// here's an array in which you can store the decoded documents
	var songs []*models.Song

	// passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := c.collection.Find(ctx, bson.D{{}}, findOptions)
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
	cur.Close(ctx)

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
	filter := bson.M{"id": song.ID}

	updateResult, err := c.collection.ReplaceOne(ctx, filter, song)
	if err != nil {
		return err
	}

	log.Printf("matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (c SongsRepository) Delete(ctx context.Context, id int) error {
	filter := bson.M{"id": id}

	deleteResult, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Printf("Deleted %v documents in the songs collection\n", deleteResult.DeletedCount)

	return nil
}
