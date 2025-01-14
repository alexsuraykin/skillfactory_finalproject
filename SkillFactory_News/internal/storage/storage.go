package storage

import (
	"context"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"
	feedsdb "github.com/alexsuraykin/SkillFactory_News/internal/storage/feedsDB"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Feeds interface {
	StoreFeeds(feeds []models.Feeds) error                                                       //Save news to storage
	Feeds(page int, pageSize int, title string, filter string) (feeds []models.Feeds, err error) //Get news from storge
	FeedById(id int) (*models.Feeds, error)                                                      //Get news by Id from storge
}

type Storage struct {
	Feeds      Feeds
	Log        zerolog.Logger
	PostgresDB *pgxpool.Pool
}

func NewRepository(ctx context.Context, postgresDB *pgxpool.Pool, log zerolog.Logger) *Storage {
	return &Storage{
		Feeds:      feedsdb.NewFeedsPostgres(postgresDB, log),
		Log:        log,
		PostgresDB: postgresDB,
	}
}
