package main

import (
	"github.com/romanzimoglyad/wow/internal/client"
	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/romanzimoglyad/wow/internal/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("configuration failed")
	}
	logger.InitZeroLog(cfg)

	// run client
	client := client.NewClient(cfg)
	err = client.Start()
	if err != nil {
		log.Error().Err(err).Msg("client error")
	}
}
