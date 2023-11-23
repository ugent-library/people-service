package models

import "context"

type PersonSuggestService interface {
	SuggestPeople(context.Context, string) ([]*Person, error)
}
