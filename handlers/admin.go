package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RenderAddReleasePage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_release.html", gin.H{
		"Title": "Add New Release",
	})
}

// AdminPanelHandler для страницы панель администратора
func AdminPanelHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем user_id из сессии
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			log.Printf("User ID is not found in session")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		// Получаем роль пользователя
		var accessRights int
		err := db.QueryRow("SELECT access_rights FROM role WHERE id = $1", userID).Scan(&accessRights)
		if err != nil {
			log.Printf("Error fetching user role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
			return
		}

		// Передаем информацию о роли в шаблон
		isAdmin := accessRights == 2

		// Рендерим админ-панель
		tmpl, err := template.ParseFiles("public/releases.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
			return
		}

		if err := tmpl.Execute(c.Writer, gin.H{
			"Title":   "Admin Panel",
			"IsAdmin": isAdmin, // передаем флаг для отображения кнопки
		}); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	}
}
