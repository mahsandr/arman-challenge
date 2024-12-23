package repository

import models "github.com/mahsandr/arman-challenge/internal/domain/models"

type SegmentRepository interface {
	SaveUserSegments(segments []models.UserSegment) error
	GetSegmentCount(segment string) (int64, error)
}
