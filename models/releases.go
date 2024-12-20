package models

// Release представляет модель релиза
type Release struct {
	ID       int64  `json:"id"`        // ID релиза
	Name     string `json:"name"`      // Название релиза
	ArtistID int64  `json:"artist_id"` // ID артиста (связь с таблицей artists)
}
