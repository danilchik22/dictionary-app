package migrations

import (
	"os"

	"gorm.io/gorm"
)

func ApplySQLMigration(db *gorm.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return db.Exec(string(sqlBytes)).Error
}
