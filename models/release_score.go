package models

type ReleaseScore struct {
	ID        int64 `json:"id"`
	Score     int64 `json:"score"`
	ReleaseID int64 `json:"release_id"`
	UserID    int64 `json:"user_id"`
}
