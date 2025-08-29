package redisClient

import (
	"dictionary_app/config"
	sl "dictionary_app/utils/logger"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

var redisClient *RedisClient

func init() {
	logger := sl.GetLogger()
	db, err := strconv.Atoi(config.GetConfig().RedisConfig.Database)
	if err != nil {
		logger.Error("Error in converting from string to int")
		return
	}
	redisClient = NewRedisClient(config.GetConfig().RedisConfig.Address, config.GetConfig().RedisConfig.Password, db)
}

func NewRedisClient(address string, password string, DB int) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       DB,
		}),
	}
}

func GetRedisClient() *RedisClient {
	return redisClient
}
