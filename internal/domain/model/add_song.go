package model

type AddSong struct {
	Song        string `json:"song" binding:"required"`
	Group       string `json:"group" binding:"required"`
	ReleaseDate string `json:"releaseDate" binding:"required"`
	Link        string `json:"link" binding:"required"`
	Text        string `json:"text" binding:"required"`
}
