package models

import "time"

type Person struct {
	Active              bool         `json:"active"`
	Name                string       `json:"name,omitempty"`
	PreferredName       string       `json:"preferred_name,omitempty"`
	GivenName           string       `json:"given_name,omitempty"`
	FamilyName          string       `json:"family_name,omitempty"`
	PreferredGivenName  string       `json:"preferred_given_name,omitempty"`
	PreferredFamilyName string       `json:"preferred_family_name,omitempty"`
	HonorificPrefix     string       `json:"honorific_prefix,omitempty"`
	Email               string       `json:"email,omitempty"`
	Attributes          []Attribute  `json:"attributes,omitempty"`
	Identifiers         []Identifier `json:"identifiers,omitempty"`
}

type PersonRecord struct {
	Person
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Identifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type Attribute struct {
	Scope string `json:"scope,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
