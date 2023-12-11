package models

import "context"

type OrganizationSuggestService interface {
	SuggestOrganizations(context.Context, OrganizationSuggestParams) ([]*Organization, error)
	RebuildAutocompleteOrganizations(context.Context) error
}

type OrganizationSuggestParams struct {
	Query string
	Limit uint32
}

func (p OrganizationSuggestParams) MergeDefault() OrganizationSuggestParams {
	limit := p.Limit
	if limit == 0 {
		limit = 20
	}
	return OrganizationSuggestParams{
		Query: p.Query,
		Limit: limit,
	}
}
