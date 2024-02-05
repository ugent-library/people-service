package models

import "time"

type Person struct {
	Name                string       `json:"name,omitempty"`
	PreferredName       string       `json:"preferredName,omitempty"`
	GivenName           string       `json:"givenName,omitempty"`
	FamilyName          string       `json:"familyName,omitempty"`
	PreferredGivenName  string       `json:"preferredGivenName,omitempty"`
	PreferredFamilyName string       `json:"preferredFamilyName,omitempty"`
	HonorificPrefix     string       `json:"honorificPrefix,omitempty"`
	Email               string       `json:"email,omitempty"`
	Active              bool         `json:"active"`
	Username            string       `json:"username,omitempty"`
	Attributes          []Attribute  `json:"attributes,omitempty"`
	Identifiers         []Identifier `json:"identifiers,omitempty"`
}

type PersonRecord struct {
	Person
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
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
