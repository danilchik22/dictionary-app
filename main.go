package main

import (
	auth "dictionary_app/auth"
	conf "dictionary_app/config"
	"dictionary_app/internal/handler"
	middleware "dictionary_app/middleware"
	mg "dictionary_app/migrations"
	"dictionary_app/storage"
	"log"
	"time"

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
	server.SetTrustedProxies([]string{"192.168.3.60"})
	server.Use(middleware.LoggerMiddleware())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://81.177.48.223:5173", "http://localhost:5173", "http://192.168.3.60:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))

	// Ensure preflight requests are handled globally
	server.OPTIONS("/*path", func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH,HEAD")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(204)
	})

	public := server.Group("/")
	{
		// Preflight for public endpoints
		public.OPTIONS("/login", func(c *gin.Context) { c.Status(204) })
		public.OPTIONS("/new_user", func(c *gin.Context) { c.Status(204) })
		public.OPTIONS("/refresh", func(c *gin.Context) { c.Status(204) })

		public.POST("/new_user", func(ctx *gin.Context) {
			auth.CreateNewUser(ctx)
		})
		public.POST("/login", auth.UserLogin)
		public.POST("/refresh", auth.Refresh)
	}

	authorized := server.Group("/api")
	authorized.Use(middleware.TokenMiddleware())
	{
		// Preflight for API endpoints
		authorized.OPTIONS("/*path", func(c *gin.Context) { c.Status(204) })

		authorized.POST("/search", handler.SearchWord)
		authorized.GET("/total_words", handler.TotalWords)
		authorized.POST("/new_word", handler.NewWord)
	}

	err := server.Run("0.0.0.0:" + conf.GetConfig().HttpConfig.Port)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
