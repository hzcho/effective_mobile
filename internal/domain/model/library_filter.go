package model

type LibraryFilter struct {
	Group   string `json:"group,omitempty"`
	Song    string `json:"song,omitempty"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}
