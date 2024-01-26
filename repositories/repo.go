package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/projects-service/sqlc"
)

type Repo struct {
	client *sqlc.Queries
	config Config
}

type Config struct {
	Conn   string
	Secret []byte
}

func New(c Config) (*Repo, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, c.Conn)
	if err != nil {
		return nil, err
	}

	client := sqlc.New(pool)

	return &Repo{
		config: c,
		client: client,
	}, nil
}
