package models

type Release struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ArtistID int64  `json:"artist_id"`
}
