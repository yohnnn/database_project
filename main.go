package main

import (
	"log"
	"rzt/config"
	"rzt/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация роутера
	router := gin.Default()

	router.LoadHTMLGlob("public/*.html")

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	// Настройка CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Настроить разрешённые домены для production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// Статические файлы
	router.Static("/public", "./public")

	// Маршруты
	router.GET("/releases", handlers.RenderReleasesPage)       // Отображение страницы релизов
	router.GET("/releases-data", handlers.GetReleasesData(db)) // API для получения данных релизов
	router.POST("/register", handlers.RegisterHandler(db))     // Регистрация
	router.POST("/login", handlers.LoginHandler(db))           // Логин
	router.GET("/rate-release", handlers.RenderRatePage)       // Страница оценки
	router.POST("/rate", handlers.RateReleaseHandler(db))      // Обработчик для записи оценки
	router.GET("/account", handlers.AccountHandler(db))
	router.GET("/release/:release_id", handlers.RenderReleaseDetails(db))
	router.GET("/admin-panel", handlers.AdminPanelHandler(db))
	router.POST("/add-release", handlers.AddReleaseHandler(db))
	router.GET("/add-release", handlers.RenderAddReleasePage) // Рендер страницы добавления релиза

	// Запуск сервера
	log.Println("Server is running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
