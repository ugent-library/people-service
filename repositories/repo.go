package repositories

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/people-service/db"
	"github.com/ugent-library/people-service/models"
)

const PersonIdentifierType = "person"

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

	// gather existing related identifiers and people

	var existingIdentifiers []db.PeopleIdentifier
	var existingPeople []int64

	for _, id := range p.Identifiers {
		recs, err := queries.GetPersonIdentifiers(ctx, db.GetPersonIdentifiersParams(id))
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		for _, rec := range recs {
			if !slices.Contains(existingIdentifiers, rec) {
				existingIdentifiers = append(existingIdentifiers, rec)
			}
			if !slices.Contains(existingPeople, rec.PersonID) {
				existingPeople = append(existingPeople, rec.PersonID)
			}
		}
	}

	// create

	if len(existingPeople) == 0 {
		personID, err := queries.CreatePerson(ctx, db.CreatePersonParams{
			Active:              p.Active,
			Name:                p.Name,
			PreferredName:       pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
			GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
			FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
			PreferredGivenName:  pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
			PreferredFamilyName: pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
			HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
			Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
			Attributes:          p.Attributes,
		})
		if err != nil {
			return err
		}

		err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
			PersonID: personID,
			Type:     PersonIdentifierType,
			Value:    uuid.NewString(),
		})
		if err != nil {
			return err
		}

		for _, id := range p.Identifiers {
			err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
				PersonID: personID,
				Type:     id.Type,
				Value:    id.Value,
			})
			if err != nil {
				return err
			}
		}

		return tx.Commit(ctx)
	}

	// or update and merge if necessary

	personID := existingPeople[0]

	err = queries.UpdatePerson(ctx, db.UpdatePersonParams{
		ID:                  personID,
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

	for _, id := range existingIdentifiers {
		if id.Type == PersonIdentifierType && id.PersonID != personID {
			err = queries.TransferPersonIdentifier(ctx, db.TransferPersonIdentifierParams{
				PersonID: personID,
				Type:     id.Type,
				Value:    id.Value,
			})
			if err != nil {
				return err
			}
		}
		if id.Type != PersonIdentifierType {
			err = queries.DeletePersonIdentifier(ctx, db.DeletePersonIdentifierParams{
				Type:  id.Type,
				Value: id.Value,
			})
			if err != nil {
				return err
			}
		}

	}

	for _, id := range existingPeople {
		if id != personID {
			if err = queries.DeletePerson(ctx, id); err != nil {
				return err
			}
		}
	}

	for _, id := range p.Identifiers {
		err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
			PersonID: personID,
			Type:     id.Type,
			Value:    id.Value,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
