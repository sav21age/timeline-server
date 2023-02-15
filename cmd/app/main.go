package main

import (
	"os"
	"time"
	"timeline/config"
	"timeline/internal/app"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	// "github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = time.RFC1123
	zerolog.MessageFieldName = "msg"
	// zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// log.Logger = log.With().Stack().Caller().Logger()
	log.Logger = log.With().Caller().Logger()

	config, err := config.NewConfig()
	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}

	app.Run(config)
}
