package repository

import (
	"song_lib/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Repositories struct {
	repository.Song
}

func NewRepositories(pool *pgxpool.Pool, log *logrus.Logger) *Repositories {
	return &Repositories{
		Song: NewSong(pool, log),
	}
}
