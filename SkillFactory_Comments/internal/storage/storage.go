package storage

import (
	"context"

	"github.com/alexsuraykin/SkillFactory_Comments/internal/storage/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Storage postgresql datastore wrapper
type Storage struct {
	l       zerolog.Logger
	pg      *pgxpool.Pool
	Queries *queries.Queries
}

// NewStorage create and init new storage instance
func NewStorage(ctx context.Context, pgConnString string, log zerolog.Logger) (*Storage, error) {
	pgConn, err := pgxpool.New(ctx, pgConnString)
	if err != nil {
		log.Error().Err(err).Msg("failed to init postgres connection")
		return nil, err
	}

	if err = pgConn.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to connect to postgres db")
		pgConn.Close()
		return nil, err
	}

	if err := migration(pgConnString); err != nil {
		log.Error().Err(err).Msg("failed to init migrations")
		return nil, err
	}

	hdl := &Storage{
		l:  log,
		pg: pgConn,
	}

	hdl.Queries = queries.New(pgConn)

	return hdl, nil
}

func (hdl *Storage) StopPG() {
	if hdl.pg != nil {
		hdl.l.Info().Msg("closing PostgreSQL connection pool")
		hdl.pg.Close()
	}
}
