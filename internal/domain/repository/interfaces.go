package repository

import (
	"context"

	models "github.com/mahsandr/arman-challenge/internal/domain/models"
)

type SegmentRepository interface {
	SaveUserSegments(ctx context.Context, segments []*models.UserSegment) error
	GetSegmentUsersCount(ctx context.Context, segment string) (uint32, error)
}
