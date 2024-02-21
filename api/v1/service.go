package api

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/ugent-library/people-service/indexes"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/repositories"
)

type Service struct {
	repo  *repositories.Repo
	index *indexes.Index
}

func NewService(repo *repositories.Repo, index *indexes.Index) *Service {
	return &Service{
		repo:  repo,
		index: index,
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

func (s *Service) GetPerson(ctx context.Context, req *GetPersonRequest) (GetPersonRes, error) {
	p, err := s.repo.GetPerson(ctx, models.Identifier(req.Identifier))
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

	res := personRecordToAPI(p)

	return &res, nil
}

func (s *Service) SearchPeople(ctx context.Context, req *SearchPeopleRequest) (*PersonHits, error) {
	hits, err := s.index.SearchPeople(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	res := &PersonHits{Hits: make([]PersonRecord, len(hits))}
	for i, p := range hits {
		res.Hits[i] = personRecordToAPI(p)
	}

	return res, nil
}

func (s *Service) AddPerson(ctx context.Context, req *AddPersonRequest) error {
	p := req.Person

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

func (s *Service) AddOrganization(ctx context.Context, req *AddOrganizationRequest) error {
	return nil
}

func personRecordToAPI(p *models.PersonRecord) PersonRecord {
	attributes := make([]Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		attributes[i] = Attribute(attr)
	}
	identifiers := make([]Identifier, len(p.Identifiers))
	for i, id := range p.Identifiers {
		identifiers[i] = Identifier(id)
	}

	return PersonRecord{
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
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}
