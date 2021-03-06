package store

import (
	"context"
	"example/hello/project/internal/models"
)

type Store interface {
	Connect(uri string) error
	Close() error

	Songs() SongsRepository
	Artists() ArtistsRepository
}

type SongsRepository interface {
	GetArtist(ctx context.Context, song *models.Song) (*models.Artist, error)
	Create(ctx context.Context, song *models.Song) error
	All(ctx context.Context, filter *models.Filter) ([]*models.Song, error)
	ByID(ctx context.Context, id int) (*models.Song, error)
	Update(ctx context.Context, song *models.Song) error
	Delete(ctx context.Context, id int) error
}

type ArtistsRepository interface {
	Create(ctx context.Context, song *models.Artist) error
	All(ctx context.Context) ([]*models.Artist, error)
	ByID(ctx context.Context, id int) (*models.Artist, error)
	Update(ctx context.Context, artist *models.Artist) error
	Delete(ctx context.Context, id int) error
}
