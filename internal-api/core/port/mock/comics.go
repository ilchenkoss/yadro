// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal-api/core/port/comics.go

// Package mock_port is a generated GoMock package.
package mock

import (
	domain "myapp/internal-api/core/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockComicsRepository is a mock of ComicsRepository interface.
type MockComicsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockComicsRepositoryMockRecorder
}

// MockComicsRepositoryMockRecorder is the mock recorder for MockComicsRepository.
type MockComicsRepositoryMockRecorder struct {
	mock *MockComicsRepository
}

// NewMockComicsRepository creates a new mock instance.
func NewMockComicsRepository(ctrl *gomock.Controller) *MockComicsRepository {
	mock := &MockComicsRepository{ctrl: ctrl}
	mock.recorder = &MockComicsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComicsRepository) EXPECT() *MockComicsRepositoryMockRecorder {
	return m.recorder
}

// GetComicsByID mocks base method.
func (m *MockComicsRepository) GetComicsByID(ID int) (*domain.Comics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComicsByID", ID)
	ret0, _ := ret[0].(*domain.Comics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComicsByID indicates an expected call of GetComicsByID.
func (mr *MockComicsRepositoryMockRecorder) GetComicsByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComicsByID", reflect.TypeOf((*MockComicsRepository)(nil).GetComicsByID), ID)
}

// GetCountComics mocks base method.
func (m *MockComicsRepository) GetCountComics() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountComics")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountComics indicates an expected call of GetCountComics.
func (mr *MockComicsRepositoryMockRecorder) GetCountComics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountComics", reflect.TypeOf((*MockComicsRepository)(nil).GetCountComics))
}

// GetMaxID mocks base method.
func (m *MockComicsRepository) GetMaxID() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaxID")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMaxID indicates an expected call of GetMaxID.
func (mr *MockComicsRepositoryMockRecorder) GetMaxID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaxID", reflect.TypeOf((*MockComicsRepository)(nil).GetMaxID))
}

// GetMissedIDs mocks base method.
func (m *MockComicsRepository) GetMissedIDs() (map[int]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMissedIDs")
	ret0, _ := ret[0].(map[int]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMissedIDs indicates an expected call of GetMissedIDs.
func (mr *MockComicsRepositoryMockRecorder) GetMissedIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMissedIDs", reflect.TypeOf((*MockComicsRepository)(nil).GetMissedIDs))
}

// InsertComics mocks base method.
func (m *MockComicsRepository) InsertComics(arg0 *[]domain.Comics) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertComics", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertComics indicates an expected call of InsertComics.
func (mr *MockComicsRepositoryMockRecorder) InsertComics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertComics", reflect.TypeOf((*MockComicsRepository)(nil).InsertComics), arg0)
}

// UpdateComicsDescriptionByID mocks base method.
func (m *MockComicsRepository) UpdateComicsDescriptionByID(ID, description string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComicsDescriptionByID", ID, description)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateComicsDescriptionByID indicates an expected call of UpdateComicsDescriptionByID.
func (mr *MockComicsRepositoryMockRecorder) UpdateComicsDescriptionByID(ID, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComicsDescriptionByID", reflect.TypeOf((*MockComicsRepository)(nil).UpdateComicsDescriptionByID), ID, description)
}
