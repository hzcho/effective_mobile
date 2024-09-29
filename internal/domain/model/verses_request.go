package model

type VersesRequest struct {
	SongID  uint64 `json:"song_id" validate:"required"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}
