package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type (
	Config struct {
		Mongo  `yaml:"mongo"`
		JWT    `yaml:"jwt"`
		Email  `yaml:"email"`
		Server `yaml:"server"`
		Client `yaml:"client"`
		Token  `yaml:"token"`
		Cookie `yaml:"cookie"`
		Cache  `yaml:"cache"`
		Cipher `yaml:"hash"`
		Log    `yaml:"log"`
	}

	Mongo struct {
		User     string `env-required:"true" env:"MONGO_USER"`
		Password string `env-required:"true" env:"MONGO_PASSWORD"`
		URL      string `env:"MONGO_URL" env-default:"mongodb://localhost:27017"`
		DBName   string `env-required:"true" yaml:"db_name"`
	}

	JWT struct {
		AccessSecret string `env-required:"true" env:"JWT_ACCESS_SECRET"`
	}

	Email struct {
		From     string `env-required:"true" env:"EMAIL_FROM"`
		Password string `env-required:"true" env:"EMAIL_PASSWORD"`
		SMTP
	}

	SMTP struct {
		Host string `env:"EMAIL_SMTP_HOST" env-default:"localhost"`
		Port string `env:"EMAIL_SMTP_PORT" env-default:"1025"`
	}

	Server struct {
		Host           string        `env:"SERVER_HOST" env-default:"http://localhost"`
		Port           string        `env:"SERVER_PORT" env-default:"8000"`
		MaxHeaderBytes int           `yaml:"max_header_bytes" env-default:"1"`
		ReadTimeout    time.Duration `yaml:"read_timeout" env-default:"10s"`
		WriteTimeout   time.Duration `yaml:"write_timeout" env-default:"10s"`
	}

	Client struct {
		Url string `env:"CLIENT_URL" env-default:"http://localhost:3000"`
	}

	Token struct {
		AccessTTL  time.Duration `yaml:"access_ttl" env-default:"15m"`
		RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"720h"`
	}

	Cache struct {
		Backend  string        `yaml:"backend" env-default:"dummy"`
		CacheTTL time.Duration `yaml:"cache_ttl" env-default:"30m"`
	}

	Cookie struct {
		Secure bool `yaml:"secure" env-default:"true"`
	}

	Cipher struct {
		Salt string `env-required:"true" env:"SALT"`
	}

	Log struct {
		Level string `yaml:"level" env-default:"error"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// err := godotenv.Load("../.env")

	_, mainFilePath, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("unable to get self filename path")
	}

	mainDirPath := filepath.Dir(mainFilePath)

	envFilePath := filepath.Join(mainDirPath, "../../.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("env load error: %w", err)
	}

	cfg := &Config{}

	// err = cleanenv.ReadConfig("../config/config.yaml", cfg)

	yamlFilePath := filepath.Join(mainDirPath, "../../config/config.yaml")
	err = cleanenv.ReadConfig(yamlFilePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	// cfgJSON, err := json.MarshalIndent(cfg, "", "  ")
	// if err != nil {
	// 	log.Debug().Err(err).Msg("")
	// }
	// fmt.Printf("%+v\n", string(cfgJSON))

	log.Debug().Msg(fmt.Sprintf("%+v", *cfg))
	// fmt.Printf("%+v\n", cfg)

	return cfg, nil
}
