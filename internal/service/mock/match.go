// Code generated by MockGen. DO NOT EDIT.
// Source: match.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	domain "timeline/internal/domain"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMatchInterface is a mock of MatchInterface interface.
type MockMatchInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMatchInterfaceMockRecorder
}

// MockMatchInterfaceMockRecorder is the mock recorder for MockMatchInterface.
type MockMatchInterfaceMockRecorder struct {
	mock *MockMatchInterface
}

// NewMockMatchInterface creates a new mock instance.
func NewMockMatchInterface(ctrl *gomock.Controller) *MockMatchInterface {
	mock := &MockMatchInterface{ctrl: ctrl}
	mock.recorder = &MockMatchInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMatchInterface) EXPECT() *MockMatchInterfaceMockRecorder {
	return m.recorder
}

// GetMatchById mocks base method.
func (m *MockMatchInterface) GetMatchById(ctx context.Context, matchId int) (domain.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchById", ctx, matchId)
	ret0, _ := ret[0].(domain.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchById indicates an expected call of GetMatchById.
func (mr *MockMatchInterfaceMockRecorder) GetMatchById(ctx, matchId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchById", reflect.TypeOf((*MockMatchInterface)(nil).GetMatchById), ctx, matchId)
}

// GetMatches mocks base method.
func (m *MockMatchInterface) GetMatches(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetMatchesQueryParams) ([]domain.MatchDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatches", ctx, competitionId, seasonId, queryParams)
	ret0, _ := ret[0].([]domain.MatchDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatches indicates an expected call of GetMatches.
func (mr *MockMatchInterfaceMockRecorder) GetMatches(ctx, competitionId, seasonId, queryParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatches", reflect.TypeOf((*MockMatchInterface)(nil).GetMatches), ctx, competitionId, seasonId, queryParams)
}

// GetMatchesDates mocks base method.
func (m *MockMatchInterface) GetMatchesDates(ctx *gin.Context, competitionId, seasonId int, queryParams domain.GetDatesQueryParams) ([]primitive.DateTime, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchesDates", ctx, competitionId, seasonId, queryParams)
	ret0, _ := ret[0].([]primitive.DateTime)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMatchesDates indicates an expected call of GetMatchesDates.
func (mr *MockMatchInterfaceMockRecorder) GetMatchesDates(ctx, competitionId, seasonId, queryParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchesDates", reflect.TypeOf((*MockMatchInterface)(nil).GetMatchesDates), ctx, competitionId, seasonId, queryParams)
}