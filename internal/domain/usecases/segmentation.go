package usecases

import (
	"context"

	"github.com/mahsandr/arman-challenge/internal/domain/models"
)

type SegmentService interface {
	AddUserSegment(ctx context.Context, segment *models.UserSegment) error
}
