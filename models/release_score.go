package models

// ReleaseScore представляет структуру для хранения оценки релиза
type ReleaseScore struct {
	ID        int64 `json:"id"`
	Score     int64 `json:"score"`      // Оценка
	ReleaseID int64 `json:"release_id"` // Ссылка на релиз
	UserID    int64 `json:"user_id"`    // Ссылка на пользователя
}
