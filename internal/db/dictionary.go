package db

import (
	"dictionary_app/internal/models"
	sl "dictionary_app/utils/logger"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type DictionaryRepository struct {
	db *gorm.DB
}

func NewDictionaryRepository(db *gorm.DB) *DictionaryRepository {
	return &DictionaryRepository{db: db}
}

func (r *DictionaryRepository) Search(query string, limit int) ([]models.Dictionary, error) {
	var results []models.Dictionary
	nanoseconds := time.Now().UnixNano()
	err := r.db.Debug().Raw(
		`SELECT id, word, definition FROM dictionary 
	     WHERE lower(word) LIKE lower($1)
	     LIMIT $2`,
		query+"%", limit,
	).Scan(&results).Error
	sl.GetLogger().Info("Time searching:" + strconv.FormatInt(time.Now().UnixNano()-nanoseconds, 10))
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *DictionaryRepository) Total() (int64, error) {
	var count int64
	err := r.db.Table("dictionary").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DictionaryRepository) AddNewWord(word string, definition string) (id int, err error) {
	var newWord = models.Dictionary{
		Word:       word,
		Definition: definition,
	}
	result := r.db.Table("dictionary").Create(&newWord)
	if result.Error != nil {
		return 0, result.Error
	}
	return newWord.ID, nil
}
