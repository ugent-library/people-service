package models

import (
	"time"
)

type Person struct {
	Name                string       `json:"name"`
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
	Identifiers         []Identifier `json:"identifiers"`
}

type PersonRecord struct {
	Name                string       `json:"name"`
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
	Identifiers         []Identifier `json:"identifiers"`
	CreatedAt           time.Time    `json:"createdAt,omitempty"`
	UpdatedAt           time.Time    `json:"updatedAt,omitempty"`
}

type Organization struct {
	Name             string            `json:"name"`
	Identifiers      []Identifier      `json:"identifiers"`
	ParentIdentifier Maybe[Identifier] `json:"parentIdentifier"`
}

type OrganizationRecord struct {
	Name        string        `json:"name"`
	Identifiers []Identifier  `json:"identifiers"`
	Parent      *Organization `json:"parent"`
	CreatedAt   time.Time     `json:"createdAt,omitempty"`
	UpdatedAt   time.Time     `json:"updatedAt,omitempty"`
}

type Maybe[T any] struct {
	Present bool
	Value   T
}

type Identifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func (i Identifier) String() string {
	return i.Type + ":" + i.Value
}

type Attribute struct {
	Scope string `json:"scope,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
