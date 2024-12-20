package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReleaseWithArtist структура для представления релиза и его автора
type ReleaseWithArtist struct {
	ID           int64   `json:"id"`
	ReleaseName  string  `json:"release_name"`
	ArtistName   string  `json:"artist_name"`
	AverageScore float64 `json:"average_score"`
}

// RenderReleasesPage рендерит HTML-страницу для релизов
func RenderReleasesPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("public/releases.html")
	if err != nil {
		// Логируем ошибку загрузки шаблона
		log.Printf("Error loading template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	// Выполняем рендеринг шаблона
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		// Логируем ошибку рендеринга
		log.Printf("Error rendering template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
		return
	}
}

// GetReleasesData возвращает данные релизов в JSON-формате
func GetReleasesData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
            SELECT 
                releases.id,
                releases.name AS release_name,
                artists.name AS artist_name,
                COALESCE(AVG(release_score.score), 0) AS average_score
            FROM releases
            LEFT JOIN artists ON releases.artist_id = artists.id
            LEFT JOIN release_score ON releases.id = release_score.release_id
            GROUP BY releases.id, artists.name
        `

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch releases"})
			return
		}
		defer rows.Close()

		var releases []ReleaseWithArtist
		for rows.Next() {
			var release ReleaseWithArtist
			if err := rows.Scan(&release.ID, &release.ReleaseName, &release.ArtistName, &release.AverageScore); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan release data"})
				return
			}
			releases = append(releases, release)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while reading rows"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"releases": releases})
	}
}
