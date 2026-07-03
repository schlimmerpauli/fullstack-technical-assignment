package products

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, query ListQuery) (ListResult, error) {
	// Business logic belongs here.
	return s.repository.List(ctx, query)
}
