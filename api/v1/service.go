package api

import (
	"context"

	"github.com/go-faster/errors"
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

func (s *Service) GetPerson(ctx context.Context, id *Identifier) (GetPersonRes, error) {
	p, err := s.repo.GetPerson(ctx, models.Identifier(*id))
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, &ErrorStatusCode{
			StatusCode: 404,
			Response: Error{
				Code:    404,
				Message: "Person not found",
			},
		}
	}
	if err != nil {
		return nil, err
	}

	attributes := make([]Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = Attribute(attr)
	}
	identifiers := make([]Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = Identifier(id)
	}

	return &PersonRecord{
		Name:                p.Name,
		PreferredName:       OptString{Set: p.PreferredName != "", Value: p.PreferredName},
		GivenName:           OptString{Set: p.GivenName != "", Value: p.GivenName},
		PreferredGivenName:  OptString{Set: p.PreferredGivenName != "", Value: p.PreferredGivenName},
		FamilyName:          OptString{Set: p.FamilyName != "", Value: p.FamilyName},
		PreferredFamilyName: OptString{Set: p.PreferredFamilyName != "", Value: p.PreferredFamilyName},
		HonorificPrefix:     OptString{Set: p.HonorificPrefix != "", Value: p.HonorificPrefix},
		Email:               OptString{Set: p.Email != "", Value: p.Email},
		Username:            OptString{Set: p.Username != "", Value: p.Username},
		Active:              p.Active,
		Attributes:          attributes,
		Identifiers:         identifiers,
	}, nil
}

func (s *Service) AddPerson(ctx context.Context, p *Person) error {
	attributes := make([]models.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = models.Attribute(attr)
	}
	identifiers := make([]models.Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = models.Identifier(id)
	}

	return s.repo.AddPerson(ctx, &models.Person{
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
	})
}
