package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ugent-library/people-service/models"
)

type Conn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

type personRow struct {
	ID                  int64
	Name                string
	PreferredName       pgtype.Text
	GivenName           pgtype.Text
	FamilyName          pgtype.Text
	PreferredGivenName  pgtype.Text
	PreferredFamilyName pgtype.Text
	HonorificPrefix     pgtype.Text
	Email               pgtype.Text
	Active              bool
	Username            pgtype.Text
	Attributes          []models.Attribute
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
	Identifiers         []models.Identifier
}

func (row personRow) toPersonRecord() *models.PersonRecord {
	return &models.PersonRecord{
		Person: models.Person{
			Name:                row.Name,
			PreferredName:       row.PreferredName.String,
			GivenName:           row.GivenName.String,
			PreferredGivenName:  row.PreferredGivenName.String,
			FamilyName:          row.FamilyName.String,
			PreferredFamilyName: row.PreferredFamilyName.String,
			HonorificPrefix:     row.HonorificPrefix.String,
			Email:               row.Email.String,
			Username:            row.Username.String,
			Active:              row.Active,
			Attributes:          row.Attributes,
			Identifiers:         row.Identifiers,
		},
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

const getPersonQuery = `
WITH identifiers AS (
	SELECT pi1.*
	FROM people_identifiers pi1
	LEFT JOIN  people_identifiers pi2 ON pi1.person_id = pi2.person_id
	WHERE pi2.type = $1 AND pi2.value = $2	
)
SELECT p.*, json_agg(json_build_object('type', i.type, 'value', i.value)) AS identifiers
FROM people p, identifiers i WHERE p.id = i.person_id
GROUP BY p.id;
`

const getAllPeopleQuery = `
SELECT p.*, json_agg(json_build_object('type', pi.type, 'value', pi.value)) AS identifiers
FROM people p
LEFT JOIN people_identifiers pi ON p.id = pi.person_id
GROUP BY p.id;
`
