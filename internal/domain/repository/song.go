package repository

import (
	"context"
	"song_lib/internal/domain/model"
)

type Song interface {
	GetSongs(ctx context.Context, filter model.LibraryFilter) ([]model.SongDetails, error)
	GetVerses(ctx context.Context, filter model.VersesRequest) ([]string, error)
	Add(ctx context.Context, song model.Song) (uint64, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, song model.Song) (model.Song, error)
}
