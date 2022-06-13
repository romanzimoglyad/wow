package main

import (
	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/romanzimoglyad/wow/internal/logger"
	"github.com/romanzimoglyad/wow/internal/server"
	"github.com/romanzimoglyad/wow/internal/storage"
	"github.com/rs/zerolog/log"
	"github.com/ztrue/shutdown"
	"syscall"
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("configuration failed")
	}
	logger.InitZeroLog(cfg)
	server := server.New(cfg, storage.New(), storage.New())
	err = server.Start()
	if err != nil {
		log.Fatal().Err(err)
	}
	shutdown.Add(func() {
		log.Info().Msg("graceful shutdown: start")
		server.Shutdown()
		log.Info().Msg("graceful shutdown: finished successful")
	})
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
