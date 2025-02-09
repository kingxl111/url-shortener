// Code generated by MockGen. DO NOT EDIT.
// Source: contracts.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	url "github.com/kingxl111/url-shortener/internal/url"
)

// MockURLRepository is a mock of URLRepository interface.
type MockURLRepository struct {
	ctrl     *gomock.Controller
	recorder *MockURLRepositoryMockRecorder
}

// MockURLRepositoryMockRecorder is the mock recorder for MockURLRepository.
type MockURLRepositoryMockRecorder struct {
	mock *MockURLRepository
}

// NewMockURLRepository creates a new mock instance.
func NewMockURLRepository(ctrl *gomock.Controller) *MockURLRepository {
	mock := &MockURLRepository{ctrl: ctrl}
	mock.recorder = &MockURLRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLRepository) EXPECT() *MockURLRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockURLRepository) Create(ctx context.Context, inputURL url.URL) (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, inputURL)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockURLRepositoryMockRecorder) Create(ctx, inputURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockURLRepository)(nil).Create), ctx, inputURL)
}

// Get mocks base method.
func (m *MockURLRepository) Get(ctx context.Context, shortenedUrl url.URL) (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, shortenedUrl)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockURLRepositoryMockRecorder) Get(ctx, shortenedUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockURLRepository)(nil).Get), ctx, shortenedUrl)
}
