package models

type Repository interface {
	PersonService
	PersonSuggestService
	OrganizationService
	OrganizationSuggestService
}
