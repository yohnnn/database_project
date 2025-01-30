package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewReleaseRequest struct {
	ReleaseName string   `json:"release_name"`
	ArtistName  string   `json:"artist_name"`
	Genres      []string `json:"genres"`
}

func AddReleaseHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req NewReleaseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var artistID int
		err := db.QueryRow(`SELECT id FROM artists WHERE name = $1`, req.ArtistName).Scan(&artistID)
		if err == sql.ErrNoRows {
			err = db.QueryRow(`INSERT INTO artists (name) VALUES ($1) RETURNING id`, req.ArtistName).Scan(&artistID)
			if err != nil {
				log.Printf("Error adding artist: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add artist"})
				return
			}
		} else if err != nil {
			log.Printf("Error fetching artist: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artist"})
			return
		}

		var releaseID int
		err = db.QueryRow(`INSERT INTO releases (name, artist_id) VALUES ($1, $2) RETURNING id`, req.ReleaseName, artistID).Scan(&releaseID)
		if err != nil {
			log.Printf("Error adding release: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add release"})
			return
		}

		for _, genre := range req.Genres {
			var genreID int
			err := db.QueryRow(`SELECT id FROM genre WHERE genre = $1`, genre).Scan(&genreID)
			if err == sql.ErrNoRows {
				err = db.QueryRow(`INSERT INTO genre (genre) VALUES ($1) RETURNING id`, genre).Scan(&genreID)
				if err != nil {
					log.Printf("Error adding genre: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add genre"})
					return
				}
			} else if err != nil {
				log.Printf("Error fetching genre: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genre"})
				return
			}

			_, err = db.Exec(`INSERT INTO release_genres (release_id, genre_id) VALUES ($1, $2)`, releaseID, genreID)
			if err != nil {
				log.Printf("Error linking release to genre: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link release to genre"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Release added successfully"})
	}
}
