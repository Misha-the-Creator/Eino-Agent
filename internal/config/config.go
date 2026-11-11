package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LLM LLMConfig `env-prefix:"LLM__"`
	App AppConfig `env-prefix:"APP__"`
	DB DBconfig `env-prefix:"DB__"`
}

type LLMConfig struct {
	APIkey    string `env:"API_KEY" env-required:"true"`
	ModelName string `env:"MODEL_NAME" env-required:"true"`
}

type DBconfig struct {
	Host    string `env:"HOST" env-required:"true"`
	Port    string `env:"PORT" env-required:"true"`
	User string `env:"USER" env-required:"true"`
	Name string `env:"NAME" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	SSLMode string `env:"SSLMODE" env-required:"true"`
}

type AppConfig struct {
	Host    string `env:"HOST" env-required:"true"`
	Port    string `env:"PORT" env-required:"true"`
	IsLocal bool   `env:"IS_LOCAL" env-default:"false"`
}

func MustLoad() *Config {
	var cfg Config
	var err error

	if _, err = os.Stat(".env"); err != nil {
		err = cleanenv.ReadEnv(&cfg)
	} else {
		err = cleanenv.ReadConfig(".env", &cfg)
	}

	if err != nil {
		fmt.Printf("Проблема при чтении конфига: %v\n", err)
	}

	return &cfg

}
