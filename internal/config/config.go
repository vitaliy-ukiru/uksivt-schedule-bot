package config

import (
	"io"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		User     string `yaml:"user" env:"PG_USER"`
		Password string `yaml:"password" env:"PG_PASSWORD"`
		Database string `yaml:"database" env:"PG_DATABASE"`
		Host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
		Port     int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
	} `yaml:"database"`
	Telegram struct {
		Token           string        `env:"BOT_TOKEN"`
		SchedulerPeriod time.Duration `yaml:"cron_period" env:"PG_PORT" env-default:"1h"`
	}
	Schedule struct {
		ApiURL string `yaml:"url"`
	} `yaml:"schedule-api"`
}

var cfg *Config

func Init(file io.Reader) (err error) {
	once.Do(func() {
		cfg = new(Config)
		if err = yaml.NewDecoder(file).Decode(cfg); err != nil {
			return
		}

	if err = cleanenv.ReadEnv(cfg); err != nil {
		return
	}

	return
}

func Get() *Config {
	if cfg == nil {
		panic("need init config")
	}
	return cfg
}
