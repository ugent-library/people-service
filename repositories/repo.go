package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/people-service/db"
	"github.com/ugent-library/people-service/models"
)

type Repo struct {
	queries *db.Queries
	config  Config
}

type Config struct {
	Conn string
}

func New(c Config) (*Repo, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, c.Conn)
	if err != nil {
		return nil, err
	}

	return &Repo{
		config:  c,
		queries: db.New(pool),
	}, nil
}

func (r *Repo) AddPerson(p *models.Person) error {
	return nil
}
