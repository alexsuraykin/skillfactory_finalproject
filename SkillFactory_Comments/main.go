package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexsuraykin/SkillFactory_Comments/config"
	"github.com/alexsuraykin/SkillFactory_Comments/internal/api"
	"github.com/alexsuraykin/SkillFactory_Comments/internal/storage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate sqlc generate

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init config")
	}

	fmt.Println(cfg)

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	logger := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	ctx := context.Background()

	dbLog := logger.With().Str("module", "storage").Logger()

	storage, err := storage.NewStorage(ctx, cfg.PgConnString, dbLog)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}

	apiLog := logger.With().Str("module", "api").Logger()
	fmt.Println(cfg.APIPort)

	APIConfig := &api.Opts{
		Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.APIPort),
		Log:     apiLog,
		Storage: storage,
	}

	server, err := api.NewAPI(APIConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start api server")
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msg("awaiting signal")

	sig := <-sigs

	log.Info().Str("signal", sig.String()).Msg("signal received")

	server.Stop()
	storage.StopPG()

	logger.Info().Msg("exiting")

}
