package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexsuraykin/SkillFactory_News/config"
	"github.com/alexsuraykin/SkillFactory_News/internal/api"
	"github.com/alexsuraykin/SkillFactory_News/internal/storage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init config")
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	logger := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	ctx := context.Background()

	dbLog := logger.With().Str("module", "storage").Logger()

	db, err := storage.NewPostgres(ctx, cfg.PgConnString, dbLog)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}

	storage := storage.NewRepository(ctx, db, dbLog)

	apiLog := logger.With().Str("module", "api").Logger()
	fmt.Println(cfg.APIPort)

	APIConfig := &api.Opts{
		Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.APIPort),
		Log:     apiLog,
		Storage: storage,
	}

	server := api.NewAPI(APIConfig)

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start api server")
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGABRT)

	logger.Info().Msg("awaiting signal")

	sig := <-sigs

	log.Info().Str("signal", sig.String()).Msg("signal received")

	server.Stop(context.Background())
	storage.StopPG()

	logger.Info().Msg("exiting")

}
