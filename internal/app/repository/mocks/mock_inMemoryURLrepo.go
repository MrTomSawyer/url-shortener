// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/repository/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/MrTomSawyer/url-shortener/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepoHandler is a mock of RepoHandler interface.
type MockRepoHandler struct {
	ctrl     *gomock.Controller
	recorder *MockRepoHandlerMockRecorder
}

// MockRepoHandlerMockRecorder is the mock recorder for MockRepoHandler.
type MockRepoHandlerMockRecorder struct {
	mock *MockRepoHandler
}

// NewMockRepoHandler creates a new mock instance.
func NewMockRepoHandler(ctrl *gomock.Controller) *MockRepoHandler {
	mock := &MockRepoHandler{ctrl: ctrl}
	mock.recorder = &MockRepoHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoHandler) EXPECT() *MockRepoHandlerMockRecorder {
	return m.recorder
}

// BatchCreate mocks base method.
func (m *MockRepoHandler) BatchCreate(data []models.TempURLBatchRequest, userID string) ([]models.BatchURLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchCreate", data, userID)
	ret0, _ := ret[0].([]models.BatchURLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchCreate indicates an expected call of BatchCreate.
func (mr *MockRepoHandlerMockRecorder) BatchCreate(data, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchCreate", reflect.TypeOf((*MockRepoHandler)(nil).BatchCreate), data, userID)
}

// Create mocks base method.
func (m *MockRepoHandler) Create(shortURL, originalURL, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", shortURL, originalURL, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepoHandlerMockRecorder) Create(shortURL, originalURL, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepoHandler)(nil).Create), shortURL, originalURL, userID)
}

// DeleteAll mocks base method.
func (m *MockRepoHandler) DeleteAll(shortURLs []string, userid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", shortURLs, userid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockRepoHandlerMockRecorder) DeleteAll(shortURLs, userid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockRepoHandler)(nil).DeleteAll), shortURLs, userid)
}

// GetAll mocks base method.
func (m *MockRepoHandler) GetAll(userid string) ([]models.URLJsonResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userid)
	ret0, _ := ret[0].([]models.URLJsonResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRepoHandlerMockRecorder) GetAll(userid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepoHandler)(nil).GetAll), userid)
}

// OriginalURL mocks base method.
func (m *MockRepoHandler) OriginalURL(shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OriginalURL", shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OriginalURL indicates an expected call of OriginalURL.
func (mr *MockRepoHandlerMockRecorder) OriginalURL(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OriginalURL", reflect.TypeOf((*MockRepoHandler)(nil).OriginalURL), shortURL)
}
