package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/people-service/db"
	"github.com/ugent-library/people-service/models"
)

type Repo struct {
	config  Config
	db      *pgxpool.Pool
	queries *db.Queries
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
		db:      pool,
		queries: db.New(pool),
	}, nil
}

func (r *Repo) AddPerson(ctx context.Context, p *models.Person) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	// TODO we only need the identifiers, not the person
	var existingRecs []db.GetPersonByIdentifierRow
	for t, vals := range p.Identifiers {
		for _, val := range vals {
			rec, err := queries.GetPersonByIdentifier(ctx, db.GetPersonByIdentifierParams{Type: t, Value: val})
			if err == nil {
				existingRecs = append(existingRecs, rec)
			} else if err != pgx.ErrNoRows {
				return err
			}
		}
	}

	if len(existingRecs) == 0 {
		id, err := queries.CreatePerson(ctx, db.CreatePersonParams{
			Active:              p.Active,
			Name:                p.Name,
			PreferredName:       pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
			GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
			FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
			PreferredGivenName:  pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
			PreferredFamilyName: pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
			HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
			Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
		})
		if err != nil {
			return err
		}

		err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
			PersonID: id,
			Type:     "person",
			Value:    uuid.NewString(),
		})
		if err != nil {
			return err
		}

		for t, vals := range p.Identifiers {
			for _, val := range vals {
				err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
					PersonID: id,
					Type:     t,
					Value:    val,
				})
				if err != nil {
					return err
				}
			}
		}

		return tx.Commit(ctx)
	}

	err = queries.UpdatePerson(ctx, db.UpdatePersonParams{
		ID:                  existingRecs[0].ID,
		Active:              p.Active,
		Name:                p.Name,
		PreferredName:       pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
		GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
		FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
		PreferredGivenName:  pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
		PreferredFamilyName: pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
		HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
		Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
		UpdatedAt:           pgtype.Timestamptz{Valid: true, Time: time.Now()},
	})
	if err != nil {
		return err
	}

	for i, rec := range existingRecs {
		var ids []struct{ Type, Value string }
		if err := json.Unmarshal(rec.Identifiers, &ids); err != nil {
			return err
		}
		for _, id := range ids {
			if id.Type == "person" && i != 0 {
				err = queries.MovePersonIdentifier(ctx, db.MovePersonIdentifierParams{
					PersonID: existingRecs[0].ID,
					Type:     id.Type,
					Value:    id.Value,
				})
				if err != nil {
					return err
				}
			}
			if id.Type != "person" {
				err = queries.DeletePersonIdentifier(ctx, db.DeletePersonIdentifierParams{
					Type:  id.Type,
					Value: id.Value,
				})
				if err != nil {
					return err
				}
			}
		}

		if i != 0 {
			if err = queries.DeletePerson(ctx, rec.ID); err != nil {
				return err
			}
		}
	}

	for t, vals := range p.Identifiers {
		for _, val := range vals {
			err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
				PersonID: existingRecs[0].ID,
				Type:     t,
				Value:    val,
			})
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
