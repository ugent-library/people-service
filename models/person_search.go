package models

import "context"

type PersonSuggestService interface {
	SuggestPeople(context.Context, PersonSuggestParams) ([]*Person, error)
	RebuildAutocompletePeople(context.Context) error
}

type PersonSuggestParams struct {
	Query  string
	Limit  uint32
	Active []bool
}

func (p PersonSuggestParams) MergeDefault() PersonSuggestParams {
	active := p.Active
	if len(active) == 0 {
		active = append(active, true, false)
	}
	limit := p.Limit
	if limit == 0 {
		limit = 20
	}
	return PersonSuggestParams{
		Query:  p.Query,
		Limit:  limit,
		Active: active,
	}
}
