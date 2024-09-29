package usecase

import (
	"song_lib/internal/domain/usecase"
	"song_lib/internal/repository"

	"github.com/sirupsen/logrus"
)

type Usecases struct {
	usecase.Song
}

func NewUsecases(repos repository.Repositories, log *logrus.Logger) *Usecases {
	return &Usecases{
		Song: NewSong(repos.Song, log),
	}
}
