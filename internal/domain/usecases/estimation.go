package usecases

import "context"

type EstimationService interface {
	GetSegmentUsersCount(ctx context.Context, segmentName string) (uint32, error)
}
