package repository

import (
	"context"
	"fmt"
	"song_lib/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Song struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewSong(pool *pgxpool.Pool, log *logrus.Logger) *Song {
	return &Song{
		pool: pool,
		log:  log,
	}
}

func (s *Song) GetSongs(ctx context.Context, filter model.LibraryFilter) ([]model.SongDetails, error) {
	log := s.log.WithField("op", "internal/repository/song/GetSongs")

	log.Debugf("Received filter: %+v", filter)

	query := "SELECT release_date, text, link FROM songs WHERE 1=1"
	var args []interface{}
	argID := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND group_name = $%d", argID)
		argID++
		args = append(args, filter.Group)
		log.Debugf("Added group filter: %s", filter.Group)
	}
	if filter.Song != "" {
		query += fmt.Sprintf(" AND song = $%d", argID)
		argID++
		args = append(args, filter.Song)
		log.Debugf("Added song filter: %s", filter.Song)
	}

	query += fmt.Sprintf(" LIMIT $%d", argID)
	argID++
	args = append(args, filter.PerPage)

	query += fmt.Sprintf(" OFFSET $%d", argID)
	args = append(args, filter.PerPage*filter.Page)

	log.Debugf("Executing query: %s with args: %+v", query, args)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var songs []model.SongDetails
	for rows.Next() {
		s := model.SongDetails{}

		err := rows.Scan(
			&s.ReleaseDate,
			&s.Text,
			&s.Link,
		)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		songs = append(songs, s)
	}

	log.Infof("Successfully retrieved %d songs", len(songs))
	return songs, nil
}

func (s *Song) GetVerses(ctx context.Context, filter model.VersesRequest) ([]string, error) {
	log := s.log.WithField("op", "internal/repository/song/GetVerses")

	log.Debugf("Received filter: %+v", filter)

	query := `
WITH split_songs AS (
    SELECT unnest(string_to_array(text, E'\n\n')) AS verse
    FROM songs
    WHERE id = $1
)
SELECT verse
FROM split_songs
ORDER BY verse
LIMIT $2 OFFSET $3;
`

	log.Debugf("Executing query: %s with args: [%d, %d, %d]", query, filter.SongID, filter.PerPage, filter.Page*filter.PerPage)

	rows, err := s.pool.Query(ctx, query, filter.SongID, filter.PerPage, filter.Page*filter.PerPage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var verses []string
	for rows.Next() {
		var verse string

		if err := rows.Scan(&verse); err != nil {
			log.Error(err)
			return nil, err
		}

		verses = append(verses, verse)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("Successfully retrieved %d verses for song ID: %d", len(verses), filter.SongID)
	return verses, nil
}

func (s *Song) Add(ctx context.Context, song model.Song) (uint64, error) {
	log := s.log.WithField("op", "internal/repository/song/Add")

	log.Debugf("Received song to add: %+v", song)

	query := "INSERT INTO songs (song, group_name, release_date, link, text) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	log.Debugf("Executing query: %s", query)

	row := s.pool.QueryRow(
		ctx,
		query,
		song.Song,
		song.Group,
		song.ReleaseDate,
		song.Link,
		song.Text,
	)

	var id uint64
	if err := row.Scan(&id); err != nil {
		log.Error(err)
		return 0, err
	}

	log.Infof("Successfully added song with ID: %d", id)
	return id, nil
}

func (s *Song) Delete(ctx context.Context, id uint64) error {
	log := s.log.WithField("op", "internal/repository/song/Delete")

	log.Infof("Attempting to delete song with ID: %d", id)

	query := "DELETE FROM songs WHERE id=$1"

	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Successfully deleted song with ID: %d", id)
	return nil
}

func (s *Song) Update(ctx context.Context, song model.Song) (model.Song, error) {
	log := s.log.WithField("op", "internal/repository/song/Update")

	log.Debugf("Received song to update: %+v", song)

	query := "UPDATE songs SET"
	var args []interface{}
	argID := 1

	if song.Group != "" {
		query += fmt.Sprintf(" group_name = $%d,", argID)
		argID++
		args = append(args, song.Group)
	}

	if song.Song != "" {
		query += fmt.Sprintf(" song = $%d,", argID)
		argID++
		args = append(args, song.Song)
	}

	if song.ReleaseDate != "" {
		query += fmt.Sprintf(" release_date = $%d,", argID)
		argID++
		args = append(args, song.ReleaseDate)
	}
	if song.Link != "" {
		query += fmt.Sprintf(" link = $%d,", argID)
		argID++
		args = append(args, song.Link)
	}
	if song.Text != "" {
		query += fmt.Sprintf(" text = $%d,", argID)
		argID++
		args = append(args, song.Text)
	}

	query = query[:len(query)-1]

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, song, group_name, release_date, link, text", argID)
	args = append(args, song.ID)

	log.Debugf("Executing query: %s with args: %+v", query, args)

	row := s.pool.QueryRow(ctx, query, args...)

	var updatedSong model.Song
	if err := row.Scan(
		&updatedSong.ID,
		&updatedSong.Song,
		&updatedSong.Group,
		&updatedSong.ReleaseDate,
		&updatedSong.Link,
		&updatedSong.Text,
	); err != nil {
		log.Error(err)
		return model.Song{}, err
	}

	log.Infof("Successfully updated song with ID: %d", updatedSong.ID)
	return updatedSong, nil
}
