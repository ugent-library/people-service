package models

import "time"

type Person struct {
	Active              bool               `json:"active"`
	Roles               []string           `json:"roles,omitempty"`
	Identifiers         []PersonIdentifier `json:"identifiers,omitempty"`
	Name                string             `json:"name,omitempty"`
	PreferredName       string             `json:"preferred_name,omitempty"`
	GivenName           string             `json:"given_name,omitempty"`
	FamilyName          string             `json:"family_name,omitempty"`
	PreferredGivenName  string             `json:"preferred_given_name,omitempty"`
	PreferredFamilyName string             `json:"preferred_family_name,omitempty"`
	HonorificPrefix     string             `json:"honorific_prefix,omitempty"`
	Email               string             `json:"email,omitempty"`
}

type PersonRecord struct {
	Person
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type PersonIdentifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
