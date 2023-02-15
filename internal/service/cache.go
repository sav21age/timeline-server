package service

import (
	"errors"
	"timeline/config"
	"timeline/pkg/cache"

	"github.com/rs/zerolog/log"
)

type CacheService struct {
	cache cache.Cache
	cfg   *config.Config
}

func NewCacheService(cfg *config.Config) *CacheService {

	switch cfg.Cache.Backend {

	case "locmem":
		return &CacheService{
			cache: cache.NewLocMem(),
			cfg:   cfg,
		}

	default:
		return &CacheService{
			cache: cache.NewDummy(),
			cfg:   cfg,
		}
	}
}

//go:generate mockgen -source=cache.go -destination=mock/cache.go
type CacheInterface interface {
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
}

//--

func (s *CacheService) Get(key string) (value interface{}, err error) {
	if s.cfg.Cache.Backend == "dummy" {
		return "", errors.New("cache backend is: dummy")
	}

	if key == "" {
		return "", errors.New("key is empty")
	}

	value, err = s.cache.Get(key)
	if err != nil {
		log.Error().Err(err).Msg("")
		return "", err
	}

	return value, nil
}

//--

func (s *CacheService) Set(key string, value interface{}) {
	if s.cfg.Cache.Backend != "dummy" {
		s.cache.Set(key, value, int64(s.cfg.Cache.CacheTTL))
	}
}
