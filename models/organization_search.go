package models

import "context"

type OrganizationSuggestService interface {
	SuggestOrganizations(context.Context, string) ([]*Organization, error)
}
