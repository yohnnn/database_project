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
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("public/*.html")

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	router.Static("/public", "./public")

	router.GET("/releases", handlers.RenderReleasesPage)
	router.GET("/releases-data", handlers.GetReleasesData(db))
	router.POST("/register", handlers.RegisterHandler(db))
	router.POST("/login", handlers.LoginHandler(db))
	router.GET("/rate-release", handlers.RenderRatePage)
	router.POST("/rate", handlers.RateReleaseHandler(db))
	router.GET("/account", handlers.AccountHandler(db))
	router.GET("/release/:release_id", handlers.RenderReleaseDetails(db))
	router.GET("/admin-panel", handlers.AdminPanelHandler(db))
	router.POST("/add-release", handlers.AddReleaseHandler(db))
	router.GET("/add-release", handlers.RenderAddReleasePage)

	log.Println("Server is running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
