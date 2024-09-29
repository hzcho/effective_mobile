package repository

import (
	"context"
	"fmt"
	"song_lib/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Song struct {
	pool *pgxpool.Pool
}

func NewSong(pool *pgxpool.Pool) *Song {
	return &Song{
		pool: pool,
	}
}

func (s *Song) GetSongs(ctx context.Context, filter model.LibraryFilter) ([]model.SongDetails, error) {
	query := "SELECT release_date, text, link FROM songs WHERE 1=1"
	var args []interface{}
	argID := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND group_name = $%d", argID)
		argID++
		args = append(args, filter.Group)
	}
	if filter.Song != "" {
		query += fmt.Sprintf(" AND song = $%d", argID)
		argID++
		args = append(args, filter.Song)
	}

	query += fmt.Sprintf(" LIMIT $%d", argID)
	argID++
	args = append(args, filter.PerPage)

	query += fmt.Sprintf(" OFFSET $%d", argID)
	args = append(args, filter.PerPage*filter.Page)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
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
			return nil, err
		}

		songs = append(songs, s)
	}

	return songs, nil
}

func (s *Song) GetVerses(ctx context.Context, filter model.VersesRequest) ([]string, error) {
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

	rows, err := s.pool.Query(ctx, query, filter.SongID, filter.PerPage, filter.Page*filter.PerPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []string
	for rows.Next() {
		var verse string

		if err := rows.Scan(&verse); err != nil {
			return nil, err
		}

		verses = append(verses, verse)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return verses, nil
}

func (s *Song) Add(ctx context.Context, song model.Song) (uint64, error) {
	query := "insert into songs (song, group_name, release_date, link, text) values($1, $2, $3, $4, $5) returning id"

	row := s.pool.QueryRow(
		ctx,
		query,
		song.Song,
		song.Group,
		song.ReleaseDate,
		song.Link,
		song.Text,
	)

	var id uint64 = 0
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Song) Delete(ctx context.Context, id uint64) error {
	query := "delete from songs where id=$1"

	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Song) Update(ctx context.Context, song model.Song) (model.Song, error) {
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
		return model.Song{}, err
	}

	return updatedSong, nil
}
