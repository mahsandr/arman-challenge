// segmentConsumer is responsible for consuming segments from a message broker,
// processing them, and saving them to the repository. It uses a wait group to
// manage concurrent processing of segments.
//
// NewSegmentConsumer creates a new instance of segmentConsumer.
//
// StartConsuming starts consuming segments from the message broker. It listens
// for segments on the segmentChannel and processes them concurrently. The
// context with cancel should be used to stop consuming and wait for all
// goroutines to finish.
//
// processSegments processes a batch of segments by unmarshaling them into
// UserSegment models and saving them to the repository. If an error occurs
// during unmarshaling or saving, it logs the error.
//
// Stop stops the segment consumer by closing the message broker and waiting for
// all goroutines to finish.
package estimation

import (
	"context"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	"github.com/mahsandr/arman-challenge/internal/domain/repository"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type segmentConsumer struct {
	repo      repository.SegmentRepository
	consumer  repository.MessageBroker
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
}

func NewSegmentConsumer(repo repository.SegmentRepository, consumer repository.MessageBroker, logger *zap.Logger) *segmentConsumer {
	return &segmentConsumer{
		repo:      repo,
		consumer:  consumer,
		logger:    logger,
		waitGroup: &sync.WaitGroup{},
	}
}

func (uc *segmentConsumer) StartConsuming(ctx context.Context) error {
	segmentChannel, err := uc.consumer.Consume(ctx)
	if err != nil {
		uc.logger.Error("error on consuming", zap.Error(err))
		return err
	}

	for {
		select {
		case <-ctx.Done():
			uc.Stop()
			return nil
		case segments := <-segmentChannel:
			uc.waitGroup.Add(1)
			go func(segments [][]byte) {
				defer uc.waitGroup.Done()
				uc.processSegments(ctx, segments)
			}(segments)
		}
	}
}
func (uc *segmentConsumer) processSegments(ctx context.Context, segments [][]byte) {
	userSegments := make([]*models.UserSegment, 0, len(segments))
	for _, segment := range segments {
		userSegment := &models.UserSegment{}
		err := json.Unmarshal(segment, userSegment)
		if err != nil {
			uc.logger.Error("error unmarshaling segment", zap.Error(err))
			continue
		}
		userSegments = append(userSegments, userSegment)
	}

	if len(userSegments) == 0 {
		return
	}
	err := uc.repo.SaveUserSegments(ctx, userSegments)
	if err != nil {
		uc.logger.Error("error saving segments", zap.Error(err))
		ctx.Done()
	}
}

func (uc *segmentConsumer) Stop() {
	uc.consumer.Close()
	uc.waitGroup.Wait()
}
