package service

import (
	"context"
	"dictionary_app/internal/db"
	"dictionary_app/internal/models"
	"dictionary_app/redisClient"
	sl "dictionary_app/utils/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type DictionaryService struct {
	repo *db.DictionaryRepository
}

type ResponseFromService struct {
	SeveralWords []models.Dictionary
	OneWord      models.Dictionary
}

func NewDictionaryService(repo *db.DictionaryRepository) *DictionaryService {
	return &DictionaryService{repo: repo}
}

func (s *DictionaryService) Search(query string, limit int, one bool) (*ResponseFromService, error) {
	logger := sl.GetLogger()
	if !one {
		results, err := s.repo.Search(query, limit)
		if err != nil {
			return &ResponseFromService{
				SeveralWords: []models.Dictionary{},
				OneWord:      models.Dictionary{},
			}, err
		}

		var oneWord models.Dictionary
		if len(results) > 0 {
			oneWord = results[0]
		}

		return &ResponseFromService{
			SeveralWords: results,
			OneWord:      oneWord,
		}, nil
	}

	redisClient := redisClient.GetRedisClient()
	cashResult, err := redisClient.Client.Get(context.Background(), query).Result()

	if err == nil {
		var cResult models.Dictionary
		json.Unmarshal([]byte(cashResult), &cResult)
		return &ResponseFromService{
			SeveralWords: []models.Dictionary{},
			OneWord:      cResult,
		}, nil
	}

	if err != nil && err != redis.Nil {
		logger.Error("error in redis client in Search Word")
		return nil, fmt.Errorf("error in redis client in Search Word: %w", err)
	}

	results, err := s.repo.Search(query, limit)
	if err != nil {
		return &ResponseFromService{
			SeveralWords: []models.Dictionary{},
			OneWord:      models.Dictionary{},
		}, err
	}

	if len(results) > 0 {
		jsonDict, err := json.Marshal(results[0])
		if err != nil {
			logger.Error("Error in setting new key in Redis: " + err.Error())
		} else {
			if err := redisClient.Client.Set(context.Background(), string(jsonDict), results[0].Word, time.Hour*240).Err(); err != nil {
				logger.Error("Error in setting new key in Redis: " + err.Error())
			}
		}
	}

	var oneWord models.Dictionary
	if len(results) > 0 {
		oneWord = results[0]
	}

	return &ResponseFromService{
		SeveralWords: results,
		OneWord:      oneWord,
	}, nil
}

func (s *DictionaryService) TotalWords() (int64, error) {
	count, err := s.repo.Total()
	if err != nil {
		return 0, err
	}
	return count, nil
}
