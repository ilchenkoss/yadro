// Code generated by MockGen. DO NOT EDIT.
// Source: internal-api/core/port/weights.go

// Package mock_port is a generated GoMock package.
package mock

import (
	domain "myapp/internal-xkcd/core/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWeightRepository is a mock of WeightRepository interface.
type MockWeightRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWeightRepositoryMockRecorder
}

// MockWeightRepositoryMockRecorder is the mock recorder for MockWeightRepository.
type MockWeightRepositoryMockRecorder struct {
	mock *MockWeightRepository
}

// NewMockWeightRepository creates a new mock instance.
func NewMockWeightRepository(ctrl *gomock.Controller) *MockWeightRepository {
	mock := &MockWeightRepository{ctrl: ctrl}
	mock.recorder = &MockWeightRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWeightRepository) EXPECT() *MockWeightRepositoryMockRecorder {
	return m.recorder
}

// GetWeightsByWords mocks base method.
func (m *MockWeightRepository) GetWeightsByWords(words map[string]float64) (*[]domain.Weights, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeightsByWords", words)
	ret0, _ := ret[0].(*[]domain.Weights)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWeightsByWords indicates an expected call of GetWeightsByWords.
func (mr *MockWeightRepositoryMockRecorder) GetWeightsByWords(words interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeightsByWords", reflect.TypeOf((*MockWeightRepository)(nil).GetWeightsByWords), words)
}

// InsertPositions mocks base method.
func (m *MockWeightRepository) InsertPositions(arg0 *[]domain.Positions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPositions", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPositions indicates an expected call of InsertPositions.
func (mr *MockWeightRepositoryMockRecorder) InsertPositions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPositions", reflect.TypeOf((*MockWeightRepository)(nil).InsertPositions), arg0)
}

// InsertWeights mocks base method.
func (m *MockWeightRepository) InsertWeights(arg0 *[]domain.Weights) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertWeights", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertWeights indicates an expected call of InsertWeights.
func (mr *MockWeightRepositoryMockRecorder) InsertWeights(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertWeights", reflect.TypeOf((*MockWeightRepository)(nil).InsertWeights), arg0)
}

// MockWeightService is a mock of WeightService interface.
type MockWeightService struct {
	ctrl     *gomock.Controller
	recorder *MockWeightServiceMockRecorder
}

// MockWeightServiceMockRecorder is the mock recorder for MockWeightService.
type MockWeightServiceMockRecorder struct {
	mock *MockWeightService
}

// NewMockWeightService creates a new mock instance.
func NewMockWeightService(ctrl *gomock.Controller) *MockWeightService {
	mock := &MockWeightService{ctrl: ctrl}
	mock.recorder = &MockWeightServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWeightService) EXPECT() *MockWeightServiceMockRecorder {
	return m.recorder
}

// FindRelevantPictures mocks base method.
func (m *MockWeightService) FindRelevantPictures(requestWeights map[string]float64, weights *[]domain.Weights) ([]domain.Comics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRelevantPictures", requestWeights, weights)
	ret0, _ := ret[0].([]domain.Comics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRelevantPictures indicates an expected call of FindRelevantPictures.
func (mr *MockWeightServiceMockRecorder) FindRelevantPictures(requestWeights, weights interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRelevantPictures", reflect.TypeOf((*MockWeightService)(nil).FindRelevantPictures), requestWeights, weights)
}

// WeightComics mocks base method.
func (m *MockWeightService) WeightComics(comics []domain.Comics) *[]domain.Weights {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WeightComics", comics)
	ret0, _ := ret[0].(*[]domain.Weights)
	return ret0
}

// WeightComics indicates an expected call of WeightComics.
func (mr *MockWeightServiceMockRecorder) WeightComics(comics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WeightComics", reflect.TypeOf((*MockWeightService)(nil).WeightComics), comics)
}

// WeightRequest mocks base method.
func (m *MockWeightService) WeightRequest(request string) map[string]float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WeightRequest", request)
	ret0, _ := ret[0].(map[string]float64)
	return ret0
}

// WeightRequest indicates an expected call of WeightRequest.
func (mr *MockWeightServiceMockRecorder) WeightRequest(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WeightRequest", reflect.TypeOf((*MockWeightService)(nil).WeightRequest), request)
}