package estimation

import (
	"context"

	"github.com/mahsandr/arman-challenge/internal/domain/repository"
	"go.uber.org/zap"
)

// Service provides methods to interact with segment data and perform estimations.
// It uses a SegmentRepository to access the data and a logger for logging errors and information.
type Service struct {
	repo   repository.SegmentRepository
	logger *zap.Logger
}

func NewService(repo repository.SegmentRepository, logger *zap.Logger) *Service {
	s := &Service{
		repo:   repo,
		logger: logger,
	}

	return s
}

// GetSegmentUsersCount retrieves the count of users in a given segment.
// It takes a context for request scoping and a segment identifier as parameters.
// It returns the count of users and an error if any occurs during the retrieval process.
func (s *Service) GetSegmentUsersCount(ctx context.Context, segment string) (uint32, error) {
	count, err := s.repo.GetSegmentUsersCount(ctx, segment)
	if err != nil {
		s.logger.Error("error on GetSegmentCount", zap.Error(err))
		return 0, err
	}
	return count, nil
}
