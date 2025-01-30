package handlers

import (
	"database/sql"
	"net/http"
	"rzt/models"
	"rzt/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		query := "SELECT id, password_hash, role_id FROM users WHERE username = $1"
		err := db.QueryRow(query, input.Username).Scan(&user.ID, &user.Password, &user.RoleID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}

		if !utils.CheckPassword(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Save()

		c.Redirect(http.StatusSeeOther, "/releases")
	}
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
			RoleID   int64  `json:"role_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil || input.Username == "" || input.Password == "" || input.Email == "" || input.RoleID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var exists bool
		queryCheck := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2)"
		if err := db.QueryRow(queryCheck, input.Username, input.Email).Scan(&exists); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already taken"})
			return
		}

		var roleExists bool
		queryRole := "SELECT EXISTS (SELECT 1 FROM role WHERE id = $1)"
		if err := db.QueryRow(queryRole, input.RoleID).Scan(&roleExists); err != nil || !roleExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}

		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		var userID int64
		query := "INSERT INTO users (username, password_hash, email, role_id) VALUES ($1, $2, $3, $4) RETURNING id"
		err = db.QueryRow(query, input.Username, hashedPassword, input.Email, input.RoleID).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", userID)
		session.Save()

		c.Redirect(http.StatusSeeOther, "/releases")
	}
}
