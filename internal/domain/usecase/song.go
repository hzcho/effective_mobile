package usecase

import (
	"context"
	"song_lib/internal/domain/model"
)

type Song interface {
	GetLib(ctx context.Context, request model.LibraryFilter) ([]model.SongDetails, error)
	GetVerses(ctx context.Context, request model.VersesRequest) (model.VersesResponse, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, song model.UpdateSong) (model.Song, error)
	Add(ctx context.Context, request model.AddSong) (uint64, error)
}
