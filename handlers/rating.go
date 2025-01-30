package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RenderRatePage(c *gin.Context) {
	tmpl, err := template.ParseFiles("public/rate.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	tmpl.Execute(c.Writer, nil)
}

func RateReleaseHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			ReleaseID int64 `json:"release_id"`
			Score     int   `json:"score"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userIDInt, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
			return
		}

		if input.Score < 1 || input.Score > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score range"})
			return
		}

		var existingScore int
		err := db.QueryRow("SELECT score FROM release_score WHERE release_id = $1 AND user_id = $2", input.ReleaseID, userIDInt).Scan(&existingScore)

		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing score"})
			return
		}

		if err == nil {
			_, err := db.Exec(`
				UPDATE release_score
				SET score = $1
				WHERE release_id = $2 AND user_id = $3
			`, input.Score, input.ReleaseID, userIDInt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update score"})
				return
			}
		} else {

			_, err := db.Exec(`
				INSERT INTO release_score (release_id, score, user_id)
				VALUES ($1, $2, $3)
			`, input.ReleaseID, input.Score, userIDInt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit score"})
				return
			}
		}

		c.Redirect(http.StatusSeeOther, "/releases")
	}
}
