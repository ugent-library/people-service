package models

import "time"

type Person struct {
	Active              bool                `json:"active"`
	Name                string              `json:"name"`
	PreferredName       string              `json:"preferred_name,omitempty"`
	GivenName           string              `json:"given_name,omitempty"`
	FamilyName          string              `json:"family_name,omitempty"`
	PreferredGivenName  string              `json:"preferred_given_name,omitempty"`
	PreferredFamilyName string              `json:"preferred_family_name,omitempty"`
	HonorificPrefix     string              `json:"honorific_prefix,omitempty"`
	Email               string              `json:"email,omitempty"`
	Roles               []string            `json:"roles,omitempty"`
	Identifiers         map[string][]string `json:"identifiers,omitempty"`
}

type PersonRecord struct {
	Person
	DateCreated time.Time `json:"date_created,omitempty"`
	DateUpdated time.Time `json:"date_updated,omitempty"`
}
