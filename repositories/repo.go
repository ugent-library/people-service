package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repo) AddPerson(ctx context.Context, p *models.Person) error {
	var existingRecs []db.GetPersonByIdentifierRow
	for t, vals := range p.Identifiers {
		for _, val := range vals {
			rec, err := r.queries.GetPersonByIdentifier(ctx, db.GetPersonByIdentifierParams{Type: t, Value: val})
			if err == nil {
				existingRecs = append(existingRecs, rec)
			} else if err != pgx.ErrNoRows {
				return err
			}
		}
	}

	if len(existingRecs) == 0 {
		params := db.CreatePersonParams{
			Active:              p.Active,
			Name:                p.Name,
			PreferredName:       pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
			GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
			FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
			PreferredGivenName:  pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
			PreferredFamilyName: pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
			HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
			Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
		}
		id, err := r.queries.CreatePerson(ctx, params)
		if err != nil {
			return err
		}

		err = r.queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
			PersonID: id,
			Type:     "person",
			Value:    uuid.NewString(),
		})
		if err != nil {
			return err
		}

		for t, vals := range p.Identifiers {
			for _, val := range vals {
				err = r.queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
					PersonID: id,
					Type:     t,
					Value:    val,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return nil
}
