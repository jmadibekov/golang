package inmemory

import (
	"context"
	"example/hello/project/internal/models"
	"fmt"
	"sync"
)

type SongsRepo struct {
	data map[int]*models.Song

	mu *sync.RWMutex
}

func (db *SongsRepo) Create(ctx context.Context, song *models.Song) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[song.ID] = song
	return nil
}

func (db *SongsRepo) All(ctx context.Context) ([]*models.Song, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	songs := make([]*models.Song, 0, len(db.data))
	for _, song := range db.data {
		songs = append(songs, song)
	}

	return songs, nil
}

func (db *SongsRepo) ByID(ctx context.Context, id int) (*models.Song, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	song, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No song with id %d", id)
	}

	return song, nil
}

func (db *SongsRepo) Update(ctx context.Context, song *models.Song) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[song.ID] = song
	return nil
}

func (db *SongsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
