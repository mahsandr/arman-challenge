package segment

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	"github.com/mahsandr/arman-challenge/internal/domain/repository"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Service struct {
	repo     repository.SegmentRepository
	producer repository.MessageBroker
	logger   *zap.Logger
}

func NewService(repo repository.SegmentRepository, producer repository.MessageBroker, logger *zap.Logger) *Service {
	s := &Service{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}

	return s
}

func (s *Service) AddUserSegment(ctx context.Context, segment *models.UserSegment) error {
	msg, err := json.Marshal(segment)
	if err != nil {
		s.logger.Error("error on marshaling", zap.Error(err))
		return err
	}
	err = s.producer.Produce(ctx, msg)
	if err != nil {
		s.logger.Error("error on producing", zap.Error(err))
		return err
	}
	return nil
}
