package storage

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"
	feedsdb "github.com/alexsuraykin/SkillFactory_News/internal/storage/feedsDB"

	"math/rand"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgres(t *testing.T) {
	conStr := "postgres://admin:admin@0.0.0.0:5432/newsAgregator?sslmode=disable"
	var log zerolog.Logger
	_, err := NewPostgres(context.Background(), conStr, log)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFeedsQueries_StoreFeeds(t *testing.T) {

	ctx := context.Background()

	// Create container for test db.
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Err(err).Msg("failed to create docker client")
	}

	containerConfig := &container.Config{
		Image: "postgres:latest",
		Env: []string{
			"POSTGRES_USER=admin",
			"POSTGRES_PASSWORD=admin",
			"POSTGRES_DB=testDB",
		},
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "5999",
				},
			},
		},
	}

	networkingConfig := &network.NetworkingConfig{}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, "my_postgres")
	if err != nil {
		log.Error().Err(err).Msg("failed to create docker container")
	}

	var options container.StartOptions
	if err := cli.ContainerStart(ctx, resp.ID, options); err != nil {
		log.Error().Err(err).Msg("failed to start docker container")
	}

	fmt.Println("PostgreSQL контейнер запущен с ID:", resp.ID)

	time.Sleep(5 * time.Second)

	// Connect to test db.
	conStr := "postgres://admin:admin@0.0.0.0:5999/testDB?sslmode=disable"
	var log zerolog.Logger
	db, err := NewPostgres(context.Background(), conStr, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
	}
	defer db.Close()

	// Up migrations
	if err := runMigrations(ctx, db); err != nil {
		log.Error().Err(err).Msg("failed to apply migrations")
	}

	repo := feedsdb.NewFeedsPostgres(db, log)

	type args struct {
		feeds []models.Feeds
	}

	testTable := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				feeds: []models.Feeds{
					{
						Title:   "testTitle",
						Content: "testContent",
						Link:    strconv.Itoa(rand.Intn(1_000_000_000)),
						PubDate: "testPubDate",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.StoreFeeds(testCase.args.feeds)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !testCase.wantErr {
				var title string
				err = db.QueryRow(context.Background(), "SELECT title FROM feeds WHERE title = $1", testCase.args.feeds[0].Title).Scan(&title)
				assert.NoError(t, err)
				assert.Equal(t, "testTitle", title)
			}
		})
	}

	var stopOpts container.StopOptions
	var rmvOpts container.RemoveOptions
	if err := cli.ContainerStop(ctx, resp.ID, stopOpts); err != nil {
		log.Error().Err(err).Msg("failed to stop docker container")
	}
	if err := cli.ContainerRemove(ctx, resp.ID, rmvOpts); err != nil {
		log.Error().Err(err).Msg("failed to remove docker container")
	}
}

func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	_, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS feeds (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		link TEXT,
		pub_date TEXT
	);`)
	if err != nil {
		return fmt.Errorf("failed to create feeds table: %w", err)
	}

	return nil
}
