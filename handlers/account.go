package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Структуры для данных
type UserRating struct {
	ReleaseName string `json:"release_name"`
	Score       int    `json:"score"`
}

type ReleaseLog struct {
	Score     int    `json:"score"`
	ReleaseID int    `json:"release_id"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// AccountHandler для страницы аккаунта
func AccountHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем user_id из сессии
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			log.Printf("User ID is not found in session")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		var username string
		var accessRights int

		// Получаем имя пользователя и роль из таблицы roles через user_id
		err := db.QueryRow(`
			SELECT u.username, r.access_rights 
			FROM users u
			JOIN role r ON u.role_id = r.id
			WHERE u.id = $1`, userID).Scan(&username, &accessRights)
		if err != nil {
			log.Printf("Error fetching username and role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username and role"})
			return
		}

		// Получаем оценки пользователя из базы данных
		query := `
            SELECT r.name, rs.score
            FROM release_score rs
            JOIN releases r ON rs.release_id = r.id
            WHERE rs.user_id = $1
        `
		rows, err := db.Query(query, userID)
		if err != nil {
			log.Printf("Error fetching user ratings: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user ratings"})
			return
		}
		defer rows.Close()
		var logs []ReleaseLog
		if accessRights == 2 { // Если пользователь администратор
			logQuery := `
        SELECT score, release_id, user_id, created_at
        FROM release_score_log
        ORDER BY created_at DESC
    `
			logRows, err := db.Query(logQuery)
			if err != nil {
				log.Printf("Error fetching release logs: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
				return
			}
			defer logRows.Close()

			for logRows.Next() {
				var logEntry ReleaseLog
				if err := logRows.Scan(&logEntry.Score, &logEntry.ReleaseID, &logEntry.UserID, &logEntry.CreatedAt); err != nil {
					log.Printf("Error scanning release logs: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan logs"})
					return
				}
				logs = append(logs, logEntry)
			}

			if err = logRows.Err(); err != nil {
				log.Printf("Log row error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading logs"})
				return
			}
		}

		var ratings []UserRating
		for rows.Next() {
			var rating UserRating
			if err := rows.Scan(&rating.ReleaseName, &rating.Score); err != nil {
				log.Printf("Error scanning user ratings: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan ratings"})
				return
			}
			ratings = append(ratings, rating)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Row error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading ratings"})
			return
		}

		// Рендеринг шаблона с оценками и ролью
		tmpl, err := template.ParseFiles("public/account.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
			return
		}

		// Передаем информацию в шаблон
		if err := tmpl.Execute(c.Writer, gin.H{
			"Title":    "My Account",
			"Username": username,
			"Ratings":  ratings,
			"IsAdmin":  accessRights == 2,
			"Logs":     logs,
		}); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	}
}
