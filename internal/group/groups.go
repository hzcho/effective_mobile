package group

import "song_lib/internal/usecase"

type Groups struct {
	Song
}

func NewGroups(usecases *usecase.Usecases) *Groups {
	return &Groups{
		Song: *NewSong(usecases.Song),
	}
}
