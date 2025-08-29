package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var once sync.Once
var ConfigVariable Config

type Config struct {
	Env               string         `yaml:"env" env:"ENV" env-default:"local"`
	HttpConfig        HTTPConfig     `yaml:"http_server" env-required:"true"`
	RedisConfig       RedisConfig    `yaml:"redis" env-required:"true"`
	DatabaseConfig    DatabaseConfig `yaml:"database" env-required:"true"`
	AuthServerAddress string         `yaml:"auth_server_address" env:"AUTH_SERVER_ADDRESS" env-required:"true"`
}

type HTTPConfig struct {
	Port        string `yaml:"port" env:"PORT" env-default:"8080"`
	Host        string `yaml:"host" env:"HOST" env-default:"localhost"`
	Timeout     string `yaml:"timeout" env:"TIMEOUT" env-default:"4"`
	TimeoutIdle string `yaml:"timeout_idle" env:"TIMEOUT_IDLE" env-default:"60"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host" env:"HOST_DATABASE" env-default:"localhost"`
	Port         int    `yaml:"port" env:"PORT_DATABASE" env-default:"5432"`
	Password     string `yaml:"password" env:"PASSWORD_DATABASE" env-required:"true"`
	Username     string `yaml:"username" env:"USERNAME_DATABASE" env-required:"true"`
	DatabaseName string `yaml:"database_name" env:"DATABASE_NAME" env-required:"true"`
}

type RedisConfig struct {
	Address  string `yaml:"address" env:"ADDRESS_REDIS" env-required:"true"`
	Password string `yaml:"password" env:"PASSWORD_REDIS" env-required:"true"`
	Database string `yaml:"database" env:"DB_REDIS" env-required:"true"`
}

func init() {
	MustLoad()
}
func MustLoad() {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatal("Path to config path is not exists")
		}
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config path does not exists: %s", configPath)
		}

		if err := cleanenv.ReadConfig(configPath, &ConfigVariable); err != nil {
			log.Fatalf("cannot read config: %s (%v)", configPath, err)
		}
	})
}

func GetConfig() *Config {
	return &ConfigVariable
}
