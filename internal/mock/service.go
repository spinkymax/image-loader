// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/spinkymax/image-loader/internal/model"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *Mockrepository) AddUser(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockrepositoryMockRecorder) AddUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*Mockrepository)(nil).AddUser), ctx, user)
}

// AuthorizeTG mocks base method.
func (m *Mockrepository) AuthorizeTG(ctx context.Context, userID int, telegramID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizeTG", ctx, userID, telegramID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthorizeTG indicates an expected call of AuthorizeTG.
func (mr *MockrepositoryMockRecorder) AuthorizeTG(ctx, userID, telegramID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeTG", reflect.TypeOf((*Mockrepository)(nil).AuthorizeTG), ctx, userID, telegramID)
}

// CheckAuth mocks base method.
func (m *Mockrepository) CheckAuth(ctx context.Context, login, password string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", ctx, login, password)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockrepositoryMockRecorder) CheckAuth(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*Mockrepository)(nil).CheckAuth), ctx, login, password)
}

// CheckTgAuth mocks base method.
func (m *Mockrepository) CheckTgAuth(ctx context.Context, tgID int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTgAuth", ctx, tgID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTgAuth indicates an expected call of CheckTgAuth.
func (mr *MockrepositoryMockRecorder) CheckTgAuth(ctx, tgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTgAuth", reflect.TypeOf((*Mockrepository)(nil).CheckTgAuth), ctx, tgID)
}

// DeleteUser mocks base method.
func (m *Mockrepository) DeleteUser(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockrepositoryMockRecorder) DeleteUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*Mockrepository)(nil).DeleteUser), ctx, id)
}

// GetAllUsers mocks base method.
func (m *Mockrepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockrepositoryMockRecorder) GetAllUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*Mockrepository)(nil).GetAllUsers), ctx)
}

// GetUser mocks base method.
func (m *Mockrepository) GetUser(ctx context.Context, id int64) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, id)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockrepositoryMockRecorder) GetUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*Mockrepository)(nil).GetUser), ctx, id)
}

// UpdateUser mocks base method.
func (m *Mockrepository) UpdateUser(ctx context.Context, modelUser model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, modelUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockrepositoryMockRecorder) UpdateUser(ctx, modelUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*Mockrepository)(nil).UpdateUser), ctx, modelUser)
}

// MockimageRepository is a mock of imageRepository interface.
type MockimageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockimageRepositoryMockRecorder
}

// MockimageRepositoryMockRecorder is the mock recorder for MockimageRepository.
type MockimageRepositoryMockRecorder struct {
	mock *MockimageRepository
}

// NewMockimageRepository creates a new mock instance.
func NewMockimageRepository(ctrl *gomock.Controller) *MockimageRepository {
	mock := &MockimageRepository{ctrl: ctrl}
	mock.recorder = &MockimageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockimageRepository) EXPECT() *MockimageRepositoryMockRecorder {
	return m.recorder
}

// AddImage mocks base method.
func (m *MockimageRepository) AddImage(ctx context.Context, modelImage model.Image) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImage", ctx, modelImage)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImage indicates an expected call of AddImage.
func (mr *MockimageRepositoryMockRecorder) AddImage(ctx, modelImage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImage", reflect.TypeOf((*MockimageRepository)(nil).AddImage), ctx, modelImage)
}

// GetImages mocks base method.
func (m *MockimageRepository) GetImages(ctx context.Context, userID int) ([]model.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImages", ctx, userID)
	ret0, _ := ret[0].([]model.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImages indicates an expected call of GetImages.
func (mr *MockimageRepositoryMockRecorder) GetImages(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImages", reflect.TypeOf((*MockimageRepository)(nil).GetImages), ctx, userID)
}

// MockfileStorage is a mock of fileStorage interface.
type MockfileStorage struct {
	ctrl     *gomock.Controller
	recorder *MockfileStorageMockRecorder
}

// MockfileStorageMockRecorder is the mock recorder for MockfileStorage.
type MockfileStorageMockRecorder struct {
	mock *MockfileStorage
}

// NewMockfileStorage creates a new mock instance.
func NewMockfileStorage(ctrl *gomock.Controller) *MockfileStorage {
	mock := &MockfileStorage{ctrl: ctrl}
	mock.recorder = &MockfileStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfileStorage) EXPECT() *MockfileStorageMockRecorder {
	return m.recorder
}

// GetObjects mocks base method.
func (m *MockfileStorage) GetObjects(ctx context.Context, images []model.Image) ([]io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjects", ctx, images)
	ret0, _ := ret[0].([]io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjects indicates an expected call of GetObjects.
func (mr *MockfileStorageMockRecorder) GetObjects(ctx, images interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjects", reflect.TypeOf((*MockfileStorage)(nil).GetObjects), ctx, images)
}

// GetUrls mocks base method.
func (m *MockfileStorage) GetUrls(ctx context.Context, images []model.Image) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUrls", ctx, images)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUrls indicates an expected call of GetUrls.
func (mr *MockfileStorageMockRecorder) GetUrls(ctx, images interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUrls", reflect.TypeOf((*MockfileStorage)(nil).GetUrls), ctx, images)
}

// PutObject mocks base method.
func (m *MockfileStorage) PutObject(ctx context.Context, image model.Image) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, image)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutObject indicates an expected call of PutObject.
func (mr *MockfileStorageMockRecorder) PutObject(ctx, image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockfileStorage)(nil).PutObject), ctx, image)
}
