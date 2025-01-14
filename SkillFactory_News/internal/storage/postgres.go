package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func NewPostgres(ctx context.Context, pgConnString string, log zerolog.Logger) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(ctx, pgConnString)
	if err != nil {
		log.Error().Err(err).Msg("failed to init postgres connection")
		log.Error().Msg("!" + pgConnString + "!")
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to connect to postgres db")
		db.Close()
		return nil, err
	}

	if err := migration(pgConnString); err != nil {
		log.Error().Err(err).Msg("failed to init migrations")
		return nil, err
	}

	return db, nil
}

func (repo *Storage) StopPG() {
	if repo.PostgresDB != nil {
		repo.Log.Info().Msg("closing PostgreSQL connection pool")
		repo.PostgresDB.Close()
	}
}
