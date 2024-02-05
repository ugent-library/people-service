package api

import (
	"context"

	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/repositories"
)

type Service struct {
	repo *repositories.Repo
}

func NewService(repo *repositories.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewError(ctx context.Context, err error) *ErrorStatusCode {
	return &ErrorStatusCode{
		StatusCode: 500,
		Response: Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}

func (s *Service) AddPerson(ctx context.Context, p *Person) error {
	return s.repo.AddPerson(ctx, toPerson(p))
}

func toPerson(p *Person) *models.Person {
	attributes := make([]models.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = models.Attribute(attr)
	}
	identifiers := make([]models.Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = models.Identifier(id)
	}

	return &models.Person{
		Name:                p.Name,
		PreferredName:       p.PreferredName.Value,
		GivenName:           p.GivenName.Value,
		PreferredGivenName:  p.PreferredGivenName.Value,
		FamilyName:          p.FamilyName.Value,
		PreferredFamilyName: p.PreferredFamilyName.Value,
		HonorificPrefix:     p.HonorificPrefix.Value,
		Email:               p.Email.Value,
		Active:              p.Active.Value,
		Username:            p.Username.Value,
		Attributes:          attributes,
		Identifiers:         identifiers,
	}
}
