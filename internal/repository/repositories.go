package repository

import (
	"song_lib/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	repository.Song
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Song: NewSong(pool),
	}
}
