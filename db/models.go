// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ugent-library/people-service/models"
)

type PeopleIdentifier struct {
	PersonID int64
	Type     string
	Value    string
}

type Person struct {
	ID                  int64
	Active              bool
	Name                string
	PreferredName       pgtype.Text
	GivenName           pgtype.Text
	FamilyName          pgtype.Text
	PreferredGivenName  pgtype.Text
	PreferredFamilyName pgtype.Text
	HonorificPrefix     pgtype.Text
	Email               pgtype.Text
	Attributes          []models.Attribute
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
}
