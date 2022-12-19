package config

import (
	"io"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
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
		AdminID         int64         `yaml:"admin-id" env:"ADMIN_ID"`
		LongPollTimeout time.Duration `yaml:"timeout" env-default:"10s"`
	}
	UKSIVT struct {
		ApiURL string `yaml:"url"`
	} `yaml:"schedule-api"`

	Scheduler struct {
		// Cron is expression for execute cron jobs
		Cron string `yaml:"cron"`
		// Range is period between time checks.
		Range        time.Duration `yaml:"range" env-default:"30m"`
		TimeLocation string        `env:"TIME_LOCATION" yaml:"time_location" env-default:"UTC"`
	} `yaml:"cron-scheduler"`
	Logger struct {
		Output struct {
			Format     string   `yaml:"format"`
			Paths      []string `yaml:"paths"`
			ErrorPaths []string `yaml:"error-paths"`
		} `yaml:"output"`

		Level zap.AtomicLevel `yaml:"level" env-default:"info"`
	} `yaml:"logger"`
}

var cfg *Config

func Init(file io.Reader) (err error) {
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
