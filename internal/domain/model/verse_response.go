package model

type VersesResponse struct {
	SongID  uint64   `json:"song_id"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Verses  []string `json:"verses"`
}
