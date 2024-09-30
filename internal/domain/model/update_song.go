package model

type UpdateSong struct {
	ID          uint64
	Song        string `json:"song,omitempty"`
	Group       string `json:"group,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Link        string `json:"link,omitempty"`
	Text        string `json:"text,omitempty"`
}

type UpdateSongSwagger struct {
	Song        string `json:"song,omitempty"`
	Group       string `json:"group,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Link        string `json:"link,omitempty"`
}
