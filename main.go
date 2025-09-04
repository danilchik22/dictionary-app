package main

import (
	auth "dictionary_app/auth"
	conf "dictionary_app/config"
	"dictionary_app/internal/handler"
	middleware "dictionary_app/middleware"
	mg "dictionary_app/migrations"
	"dictionary_app/storage"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db := storage.GetDb()
	if err := mg.ApplySQLMigration(db, "migrations/001_initial.sql"); err != nil {
		log.Fatal(err)
	}
	storage.InitNewDatabase()
	server := gin.Default()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	server.SetTrustedProxies([]string{"192.168.3.60"})
	server.Use(middleware.LoggerMiddleware())
	public := server.Group("/")
	{
		public.POST("/new_user", func(ctx *gin.Context) {
			auth.CreateNewUser(ctx)
		})
		public.POST("/login", auth.UserLogin)
		public.POST("/refresh", auth.Refresh)
	}

	authorized := server.Group("/api")
	authorized.Use(middleware.TokenMiddleware())
	{
		authorized.POST("/search", handler.SearchWord)
		authorized.GET("/total_words", handler.TotalWords)
		authorized.POST("/new_word", handler.NewWord)
	}

	err := server.Run("0.0.0.0:" + conf.GetConfig().HttpConfig.Port)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
