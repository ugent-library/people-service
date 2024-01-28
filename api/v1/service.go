package api

import (
	"context"

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

func (s *Service) AddPerson(ctx context.Context, p *Person) (*Person, error) {
	return nil, nil
}
