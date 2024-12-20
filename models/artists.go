package models

// Artist представляет модель артиста
type Artist struct {
	ID   int64  `json:"id"`   // ID артиста
	Name string `json:"name"` // Имя артиста
}
