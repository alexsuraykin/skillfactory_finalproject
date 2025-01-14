package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"skillfactory/SkillFactory_finalProject/APIGateway/config"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate ./oapigen.sh

//go:embed swagger.yaml
var swaggerFile []byte

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init config")
	}

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	logger := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	apiLog := logger.With().Str("module", "api").Logger()

	APIConfig := &api.Opts{
		Addr:        fmt.Sprintf("0.0.0.0:%v", cfg.APIPort),
		Log:         apiLog,
		SwaggerFile: swaggerFile,
	}

	server, err := api.NewAPI(APIConfig)

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

	logger.Info().Msg("exiting")
}
