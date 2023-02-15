package app

import (
	"os"
	"os/signal"
	"syscall"

	"timeline/config"

	"timeline/internal/handler"
	"timeline/internal/repository"
	"timeline/internal/service"
	"timeline/pkg/mongodb"
	"timeline/pkg/server"

	"github.com/rs/zerolog/log"
)

func Run(cfg *config.Config) {
	client, err := mongodb.NewClient(
		cfg.Mongo.URL,
		cfg.Mongo.User,
		cfg.Mongo.Password,
	)

	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
	
	db := client.Database(cfg.Mongo.DBName)

	repositories := repository.NewRepository(db)
	services := service.NewService(repositories, cfg)
	handlers := handler.NewHandler(services, cfg)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg, handlers.InitRoute()); err != nil {
			log.Error().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
