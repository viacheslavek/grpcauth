package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const defaultConfigPath = "config/local.yaml"

type Config struct {
	Env      string        `yaml:"env"`
	DB       StorageConfig `yaml:"storage"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type StorageConfig struct {
	Type     string `yaml:"type"`
	DBName   string `yaml:"db-name"`
	User     string `yaml:"user"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return mustLoadPath(configPath)
}

func mustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	if err := godotenv.Load(); err != nil {
		panic("env file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read env variables: " + err.Error())
	}

	return &cfg
}

// Priority: flag > env > default.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = defaultConfigPath
	}

	return res
}
