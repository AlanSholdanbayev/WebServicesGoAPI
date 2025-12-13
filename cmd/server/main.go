package main

import (
	"finalproject/internal/config"
	"finalproject/internal/logger"
	"finalproject/internal/server"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// используем конфиг Logger из config
	lg := logger.New(&logger.Config{
		SeqURL:    cfg.Logger.SeqURL,
		SeqAPIKey: cfg.Logger.SeqAPIKey,
	})

	srv := server.New(cfg, lg)

	if err := srv.Run(); err != nil {
		lg.Fatal().Err(err).Msg("server stopped")
	}
}
