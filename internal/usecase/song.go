package usecase

import (
	"context"
	"errors"
	"song_lib/internal/domain/model"
	"song_lib/internal/domain/repository"

	"github.com/sirupsen/logrus"
)

type Song struct {
	songRepo repository.Song
	log      *logrus.Logger
}

func NewSong(songRepo repository.Song, log *logrus.Logger) *Song {
	return &Song{
		songRepo: songRepo,
		log:      log,
	}
}

func (s *Song) GetLib(ctx context.Context, request model.LibraryFilter) ([]model.SongDetails, error) {
	log := s.log.WithField("op", "internal/usecase/building/GetLib")

	if request.PerPage <= 0 {
		request.PerPage = 10
	}
	if request.Page < 0 {
		request.PerPage = 0
	}
	songs, err := s.songRepo.GetSongs(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return songs, nil
}

func (s *Song) GetVerses(ctx context.Context, request model.VersesRequest) (model.VersesResponse, error) {
	log := s.log.WithField("op", "internal/usecase/building/GetVerses")

	if request.PerPage <= 0 {
		request.PerPage = 10
	}
	if request.Page < 0 {
		request.Page = 0
	}

	verses, err := s.songRepo.GetVerses(ctx, request)
	if err != nil {
		log.Error(err)
		return model.VersesResponse{}, err
	}

	return model.VersesResponse{
		SongID:  request.SongID,
		Page:    request.Page,
		PerPage: request.PerPage,
		Verses:  verses,
	}, nil
}

func (s *Song) Delete(ctx context.Context, id uint64) error {
	log := s.log.WithField("op", "internal/usecase/building/Delete")

	if err := s.songRepo.Delete(ctx, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *Song) Update(ctx context.Context, song model.UpdateSong) (model.Song, error) {
	log := s.log.WithField("op", "internal/usecase/building/Update")

	if song == (model.UpdateSong{}) {
		err := errors.New("No change")

		log.Error(err)
		return model.Song{}, nil
	}

	sng := model.Song{
		ID:          song.ID,
		Song:        song.Song,
		Group:       song.Group,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Song,
		Text:        song.Text,
	}

	sng, err := s.songRepo.Update(ctx, sng)
	if err != nil {
		log.Error(err)
		return model.Song{}, err
	}

	return sng, nil
}

func (s *Song) Add(ctx context.Context, request model.AddSong) (uint64, error) {
	log := s.log.WithField("op", "internal/usecase/building/Add")

	song := model.Song{
		Group:       request.Group,
		Song:        request.Song,
		ReleaseDate: request.ReleaseDate,
		Link:        request.Link,
		Text:        request.Text,
	}

	id, err := s.songRepo.Add(ctx, song)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}
