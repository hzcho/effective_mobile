package group

import (
	"song_lib/internal/usecase"

	"github.com/sirupsen/logrus"
)

type Groups struct {
	Song
}

func NewGroups(usecases *usecase.Usecases, log *logrus.Logger) *Groups {
	return &Groups{
		Song: *NewSong(usecases.Song, log),
	}
}
