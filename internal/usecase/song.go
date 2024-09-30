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
	log := s.log.WithField("op", "internal/usecase/song/GetLib")

	log.Debugf("Received request: %+v", request)

	if request.PerPage <= 0 {
		request.PerPage = 10
		log.Infof("PerPage was set to default value: %d", request.PerPage)
	}
	if request.Page < 0 {
		request.Page = 0
		log.Infof("Page was set to default value: %d", request.Page)
	}

	songs, err := s.songRepo.GetSongs(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("Successfully retrieved %d songs", len(songs))
	return songs, nil
}

func (s *Song) GetVerses(ctx context.Context, request model.VersesRequest) (model.VersesResponse, error) {
	log := s.log.WithField("op", "internal/usecase/song/GetVerses")

	log.Debugf("Received request: %+v", request)

	if request.PerPage <= 0 {
		request.PerPage = 10
		log.Infof("PerPage was set to default value: %d", request.PerPage)
	}
	if request.Page < 0 {
		request.Page = 0
		log.Infof("Page was set to default value: %d", request.Page)
	}

	verses, err := s.songRepo.GetVerses(ctx, request)
	if err != nil {
		log.Error(err)
		return model.VersesResponse{}, err
	}

	log.Infof("Successfully retrieved %d verses for SongID: %d", len(verses), request.SongID)

	return model.VersesResponse{
		SongID:  request.SongID,
		Page:    request.Page,
		PerPage: request.PerPage,
		Verses:  verses,
	}, nil
}

func (s *Song) Delete(ctx context.Context, id uint64) error {
	log := s.log.WithField("op", "internal/usecase/song/Delete")

	log.Infof("Attempting to delete song with ID: %d", id)

	if err := s.songRepo.Delete(ctx, id); err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Successfully deleted song with ID: %d", id)
	return nil
}

func (s *Song) Update(ctx context.Context, song model.UpdateSong) (model.Song, error) {
	log := s.log.WithField("op", "internal/usecase/song/Update")

	if song == (model.UpdateSong{}) {
		err := errors.New("No change")
		log.Warn(err)
		return model.Song{}, nil
	}

	sng := model.Song{
		ID:          song.ID,
		Song:        song.Song,
		Group:       song.Group,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		Text:        song.Text,
	}

	log.Infof("Updating song with ID: %d", sng.ID)

	sng, err := s.songRepo.Update(ctx, sng)
	if err != nil {
		log.Error(err)
		return model.Song{}, err
	}

	log.Infof("Successfully updated song with ID: %d", sng.ID)
	return sng, nil
}

func (s *Song) Add(ctx context.Context, request model.AddSong) (uint64, error) {
	log := s.log.WithField("op", "internal/usecase/song/Add")

	log.Debugf("Received request to add song: %+v", request)

	song := model.Song{
		Group:       request.Group,
		Song:        request.Song,
		ReleaseDate: request.ReleaseDate,
		Link:        request.Link,
		Text:        request.Text,
	}

	log.Infof("Attempting to add song: %s by group: %s", song.Song, song.Group)

	id, err := s.songRepo.Add(ctx, song)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	log.Infof("Successfully added song with ID: %d", id)
	return id, nil
}
