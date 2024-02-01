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
	identifiers := make([]models.PersonIdentifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = models.PersonIdentifier(id)
	}

	return &models.Person{
		Active:              p.Active,
		Roles:               p.Roles,
		Identifiers:         identifiers,
		Name:                p.Name,
		PreferredName:       p.PreferredName.Value,
		GivenName:           p.GivenName.Value,
		PreferredGivenName:  p.PreferredName.Value,
		FamilyName:          p.FamilyName.Value,
		PreferredFamilyName: p.PreferredFamilyName.Value,
		HonorificPrefix:     p.HonorificPrefix.Value,
		Email:               p.Email.Value,
	}
}
