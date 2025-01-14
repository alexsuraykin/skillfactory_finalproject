package feedsdb

import (
	"context"
	"fmt"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type FeedsPostgres struct {
	db *pgxpool.Pool
	l  zerolog.Logger
}

func NewFeedsPostgres(db *pgxpool.Pool, log zerolog.Logger) *FeedsPostgres {
	return &FeedsPostgres{
		db: db,
		l:  log,
	}
}

func (f *FeedsPostgres) Feeds(page int, pageSize int, title string, filter string) (feeds []models.Feeds, err error) {
	offset := (page - 1) * int(pageSize)

	rows, err := f.db.Query(context.Background(), `
		SELECT 
			id, 
			title,
			content,
			pub_date,
			link
		FROM feeds
		WHERE (title ILIKE '%' || $1 || '%') AND (content ILIKE '%' || $2 || '%')
		ORDER BY pub_date DESC
		LIMIT $3 OFFSET $4;
	`, title, filter, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Feeds
		err = rows.Scan(
			&f.Id,
			&f.Title,
			&f.Content,
			&f.PubDate,
			&f.Link,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to get feeds from storage")
			return nil, err
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}

func (f *FeedsPostgres) FeedById(id int) (*models.Feeds, error) {
	var feed models.Feeds

	err := f.db.QueryRow(context.Background(), `
		SELECT 
			id, 
			title,
			content,
			pub_date,
			link
		FROM feeds
		WHERE id = $1;
	`, id).Scan(&feed.Id, &feed.Title, &feed.Content, &feed.PubDate, &feed.Link)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no feed found with id %d", id)
		}
		return nil, fmt.Errorf("error fetching feed by id: %w", err)
	}

	return &feed, nil
}

func (f *FeedsPostgres) StoreFeeds(feeds []models.Feeds) error {
	for _, feed := range feeds {
		_, err := f.db.Exec(context.Background(), `
			INSERT INTO feeds (title, content, pub_date, link)
			VALUES ($1, $2, $3, $4)
		`, feed.Title, feed.Content, feed.PubDate, feed.Link)
		if err != nil {
			log.Error().Err(err).Msg("failed to store feeds to db")
			return err
		}
	}

	return nil
}
