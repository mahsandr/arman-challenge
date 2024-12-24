// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/domain/repository/interfaces.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/mahsandr/arman-challenge/internal/domain/models"
)

// MockSegmentRepository is a mock of SegmentRepository interface.
type MockSegmentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSegmentRepositoryMockRecorder
}

// MockSegmentRepositoryMockRecorder is the mock recorder for MockSegmentRepository.
type MockSegmentRepositoryMockRecorder struct {
	mock *MockSegmentRepository
}

// NewMockSegmentRepository creates a new mock instance.
func NewMockSegmentRepository(ctrl *gomock.Controller) *MockSegmentRepository {
	mock := &MockSegmentRepository{ctrl: ctrl}
	mock.recorder = &MockSegmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSegmentRepository) EXPECT() *MockSegmentRepositoryMockRecorder {
	return m.recorder
}

// GetSegmentUsersCount mocks base method.
func (m *MockSegmentRepository) GetSegmentUsersCount(ctx context.Context, segment string) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentUsersCount", ctx, segment)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentUsersCount indicates an expected call of GetSegmentUsersCount.
func (mr *MockSegmentRepositoryMockRecorder) GetSegmentUsersCount(ctx, segment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentUsersCount", reflect.TypeOf((*MockSegmentRepository)(nil).GetSegmentUsersCount), ctx, segment)
}

// SaveUserSegments mocks base method.
func (m *MockSegmentRepository) SaveUserSegments(ctx context.Context, segments []*models.UserSegment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserSegments", ctx, segments)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUserSegments indicates an expected call of SaveUserSegments.
func (mr *MockSegmentRepositoryMockRecorder) SaveUserSegments(ctx, segments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserSegments", reflect.TypeOf((*MockSegmentRepository)(nil).SaveUserSegments), ctx, segments)
}
