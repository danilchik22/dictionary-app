package storage

import (
	"dictionary_app/config"
	sl "dictionary_app/utils/logger"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	cfg := config.GetConfig()
	db = getDatabaseConnection(*cfg)
}
func getDatabaseConnection(cfg config.Config) *gorm.DB {
	logger := sl.GetLogger()
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.DatabaseName, cfg.DatabaseConfig.Port)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		logger.Error("Error during initialization of database", sl.Err(err))
		os.Exit(1)
	}
	return db
}

func GetDb() *gorm.DB {
	return db
}
