package sl

import (
	"dictionary_app/config"
	"log/slog"
	"os"
	"sync"
)

const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

var (
	Logger *slog.Logger
	once   sync.Once
)

func init() {
	cfg := config.GetConfig()
	InitLogger(cfg.Env)
}

func InitLogger(env string) {
	once.Do(func() {
		switch env {
		case localEnv:
			Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case devEnv:
			Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		case prodEnv:
			Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}
	})
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func GetLogger() *slog.Logger {
	return Logger
}
