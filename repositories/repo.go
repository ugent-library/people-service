package repositories

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ugent-library/people-service/db"
	"github.com/ugent-library/people-service/models"
)

const (
	IDType = "id"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repo struct {
	config  Config
	db      *pgxpool.Pool
	queries *db.Queries
}

type Config struct {
	Conn            string
	DeactivateAfter time.Duration
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

// TODO tx needed?
func (r *Repo) GetPerson(ctx context.Context, id models.Identifier) (*models.PersonRecord, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	person, err := queries.GetPerson(ctx, db.GetPersonParams(id))
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	ids, err := queries.GetPersonIdentifiers(ctx, person.ID)
	if err != nil {
		return nil, err
	}

	p := &models.PersonRecord{
		Person: models.Person{
			Name:                person.Name,
			PreferredName:       person.PreferredName.String,
			GivenName:           person.GivenName.String,
			PreferredGivenName:  person.PreferredGivenName.String,
			FamilyName:          person.FamilyName.String,
			PreferredFamilyName: person.PreferredFamilyName.String,
			HonorificPrefix:     person.HonorificPrefix.String,
			Email:               person.Email.String,
			Username:            person.Username.String,
			Active:              person.Active,
			Identifiers:         make([]models.Identifier, len(ids)),
			Attributes:          person.Attributes,
		},
		CreatedAt: person.CreatedAt.Time,
		UpdatedAt: person.UpdatedAt.Time,
	}

	for i, id := range ids {
		p.Person.Identifiers[i] = models.Identifier{Type: id.Type, Value: id.Value}
	}

	return p, tx.Commit(ctx)
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
			Name:                p.Name,
			PreferredName:       pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName},
			GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
			PreferredGivenName:  pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName},
			FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
			PreferredFamilyName: pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName},
			HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
			Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
			Active:              p.Active,
			Username:            pgtype.Text{Valid: p.Username != "", String: p.Username},
			Attributes:          p.Attributes,
		})
		if err != nil {
			return err
		}

		err = queries.CreatePersonIdentifier(ctx, db.CreatePersonIdentifierParams{
			PersonID: personID,
			Type:     IDType,
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

	preferredName := pgtype.Text{Valid: p.PreferredName != "", String: p.PreferredName}
	preferredGivenName := pgtype.Text{Valid: p.PreferredGivenName != "", String: p.PreferredGivenName}
	preferredFamilyName := pgtype.Text{Valid: p.PreferredFamilyName != "", String: p.PreferredFamilyName}

	// merge preferred names if none are given, new to old
	if !preferredName.Valid && !preferredGivenName.Valid && !preferredFamilyName.Valid {
		for _, rec := range existingPeople {
			if rec.PreferredName.String != "" || rec.PreferredGivenName.String != "" || rec.PreferredFamilyName.String != "" {
				preferredName = rec.PreferredName
				preferredFamilyName = rec.PreferredFamilyName
				preferredGivenName = rec.PreferredGivenName
				break
			}
		}
	}

	// merge attributes with non overlapping scopes, new to old
	var attributes []models.Attribute

	if len(p.Attributes) > 0 {
		attributes = append(attributes, p.Attributes...)
	}

	for _, rec := range existingPeople {
		var attrs []models.Attribute
		for _, attr := range rec.Attributes {
			if !slices.ContainsFunc(attributes, func(a models.Attribute) bool { return a.Scope == attr.Scope }) {
				attrs = append(attrs, attr)
			}
		}
		if len(attrs) > 0 {
			attributes = append(attributes, attrs...)
		}
	}

	err = queries.UpdatePerson(ctx, db.UpdatePersonParams{
		ID:                  personID,
		Name:                p.Name,
		PreferredName:       preferredName,
		GivenName:           pgtype.Text{Valid: p.GivenName != "", String: p.GivenName},
		FamilyName:          pgtype.Text{Valid: p.FamilyName != "", String: p.FamilyName},
		PreferredGivenName:  preferredGivenName,
		PreferredFamilyName: preferredFamilyName,
		HonorificPrefix:     pgtype.Text{Valid: p.HonorificPrefix != "", String: p.HonorificPrefix},
		Email:               pgtype.Text{Valid: p.Email != "", String: p.Email},
		Active:              p.Active,
		Username:            pgtype.Text{Valid: p.Username != "", String: p.Username},
		Attributes:          attributes,
	})
	if err != nil {
		return err
	}

	for _, id := range existingIdentifiers {
		if id.Type == IDType && id.PersonID != personID {
			err = queries.TransferPersonIdentifier(ctx, db.TransferPersonIdentifierParams{
				PersonID: personID,
				Type:     id.Type,
				Value:    id.Value,
			})
			if err != nil {
				return err
			}
		}
		if id.Type != IDType {
			err = queries.DeletePersonIdentifier(ctx, db.DeletePersonIdentifierParams{
				Type:  id.Type,
				Value: id.Value,
			})
			if err != nil {
				return err
			}
		}
	}

	for _, rec := range existingPeople {
		if rec.ID != personID {
			if err = queries.DeletePerson(ctx, rec.ID); err != nil {
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

func (r *Repo) DeactivatePeople(ctx context.Context) error {
	t := time.Now().Add(-r.config.DeactivateAfter)
	return r.queries.DeactivatePeople(ctx, pgtype.Timestamptz{Valid: true, Time: t})
}
