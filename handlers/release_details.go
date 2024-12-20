package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Структуры для данных релиза
type ReleaseDetails struct {
	ReleaseName   string         `json:"release_name"`
	Genres        []string       `json:"genres"`
	AverageScore  float64        `json:"average_score"`
	RatingDetails []RatingDetail `json:"user_ratings"`
}

type RatingDetail struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// RenderReleaseDetails отображает информацию о релизе
func RenderReleaseDetails(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		releaseID := c.Param("release_id")

		// Получаем информацию о релизе, жанры, среднюю оценку и пользователей с их оценками
		query := `
			SELECT 
				r.name AS release_name,
				COALESCE(AVG(rs.score), 0) AS average_score
			FROM releases r
			LEFT JOIN release_score rs ON r.id = rs.release_id
			WHERE r.id = $1
			GROUP BY r.id
		`

		var release ReleaseDetails
		err := db.QueryRow(query, releaseID).Scan(&release.ReleaseName, &release.AverageScore)
		if err != nil {
			log.Printf("Error fetching release details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch release details"})
			return
		}

		// Получаем жанры релиза
		genreQuery := `
			SELECT g.genre
			FROM genre g
			JOIN release_genres rg ON g.id = rg.genre_id
			WHERE rg.release_id = $1
		`
		genreRows, err := db.Query(genreQuery, releaseID)
		if err != nil {
			log.Printf("Error fetching genres: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
			return
		}
		defer genreRows.Close()

		var genres []string
		for genreRows.Next() {
			var genre string
			if err := genreRows.Scan(&genre); err != nil {
				log.Printf("Error scanning genre: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan genres"})
				return
			}
			genres = append(genres, genre)
		}
		release.Genres = genres

		// Получаем пользователей, которые оценили этот релиз и их оценки
		userRatingsQuery := `
			SELECT u.username, rs.score
			FROM release_score rs
			JOIN users u ON rs.user_id = u.id
			WHERE rs.release_id = $1
		`
		userRatingsRows, err := db.Query(userRatingsQuery, releaseID)
		if err != nil {
			log.Printf("Error fetching user ratings: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user ratings"})
			return
		}
		defer userRatingsRows.Close()

		var ratingDetails []RatingDetail
		for userRatingsRows.Next() {
			var rating RatingDetail
			if err := userRatingsRows.Scan(&rating.Username, &rating.Score); err != nil {
				log.Printf("Error scanning user ratings: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user ratings"})
				return
			}
			ratingDetails = append(ratingDetails, rating)
		}
		release.RatingDetails = ratingDetails

		// Рендеринг HTML-шаблона с подробной информацией о релизе
		tmpl, err := template.ParseFiles("public/release_details.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
			return
		}

		if err := tmpl.Execute(c.Writer, gin.H{
			"Release": release,
		}); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	}
}
