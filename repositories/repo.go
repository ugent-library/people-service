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

	// gather existing related people and identifiers
	var existingPeople []db.Person
	var existingIdentifiers []db.PeopleIdentifier

	for _, id := range p.Identifiers {
		person, err := queries.GetPerson(ctx, db.GetPersonParams(id))
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			continue
		}
		if !slices.ContainsFunc(existingPeople, func(p db.Person) bool { return p.ID == person.ID }) {
			identifiers, err := queries.GetPersonIdentifiers(ctx, person.ID)
			if err != nil {
				return err
			}
			existingPeople = append(existingPeople, person)
			existingIdentifiers = append(existingIdentifiers, identifiers...)
		}
	}

	slices.SortFunc(existingPeople, func(a, b db.Person) int {
		if a.UpdatedAt.Time.Before(b.UpdatedAt.Time) {
			return 1
		}
		return -1
	})

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

	personID := existingPeople[0].ID

	attrs := p.Attributes

	// TODO merge attributes with non conflicting scopes

	// TODO merge preferred names

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
		Attributes:          attrs,
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

	for _, person := range existingPeople {
		if person.ID != personID {
			if err = queries.DeletePerson(ctx, person.ID); err != nil {
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
