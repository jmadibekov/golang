package inmemory

import (
	"context"
	"example/hello/project/internal/models"
	"fmt"
	"sync"
)

type ArtistsRepo struct {
	data map[int]*models.Artist

	mu *sync.RWMutex
}

func (db *ArtistsRepo) Create(ctx context.Context, artist *models.Artist) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[artist.ID] = artist
	return nil
}

func (db *ArtistsRepo) All(ctx context.Context) ([]*models.Artist, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	artists := make([]*models.Artist, 0, len(db.data))
	for _, artist := range db.data {
		artists = append(artists, artist)
	}

	return artists, nil
}

func (db *ArtistsRepo) ByID(ctx context.Context, id int) (*models.Artist, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	artist, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No artist with id %d", id)
	}

	return artist, nil
}

func (db *ArtistsRepo) Update(ctx context.Context, artist *models.Artist) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[artist.ID] = artist
	return nil
}

func (db *ArtistsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
