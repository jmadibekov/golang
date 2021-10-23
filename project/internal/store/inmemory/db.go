package inmemory

import (
	"example/hello/project/internal/models"
	"example/hello/project/internal/store"
	"sync"
)

type DB struct {
	songsRepo   store.SongsRepository
	artistsRepo store.ArtistsRepository

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Songs() store.SongsRepository {
	if db.songsRepo == nil {
		db.songsRepo = &SongsRepo{
			data: make(map[int]*models.Song),
			mu:   new(sync.RWMutex),
		}
	}

	return db.songsRepo
}

func (db *DB) Artists() store.ArtistsRepository {
	if db.artistsRepo == nil {
		db.artistsRepo = &ArtistsRepo{
			data: make(map[int]*models.Artist),
			mu:   new(sync.RWMutex),
		}
	}

	return db.artistsRepo
}
