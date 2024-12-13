package static

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStaticUseCase struct {
	mock.Mock
}

func (m *mockStaticUseCase) GetStatic(id uuid.UUID) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *mockStaticUseCase) UploadStatic(reader io.ReadSeeker) (uuid.UUID, error) {
	args := m.Called(reader)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockStaticUseCase) GetStaticFile(uri string) (io.ReadSeeker, error) {
	args := m.Called(uri)
	return args.Get(0).(io.ReadSeeker), args.Error(1)
}

func (m *mockStaticUseCase) GetAvatar(id uuid.UUID) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestGetStatic_Success(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	staticID := uuid.New()
	expectedURI := "http://example.com/static/" + staticID.String()
	mockUC.On("GetStatic", staticID).Return(expectedURI, nil)

	result, err := grpcServer.GetStatic(context.Background(), &staticProto.Static{Id: staticID.String()})

	assert.NoError(t, err)
	assert.Equal(t, expectedURI, result.Uri)
	mockUC.AssertExpectations(t)
}

func TestGetStatic_InvalidID(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	result, err := grpcServer.GetStatic(context.Background(), &staticProto.Static{Id: "invalid-id"})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetStaticFile_Success(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	staticURI := "http://example.com/static/file"
	mockUC.On("GetStaticFile", staticURI).Return(bytes.NewReader([]byte("file content")), nil)

	stream := &mockStream{}
	stream.On("Send", &staticProto.StaticUpload{Chunk: []byte("file content")}).Return(nil)

	err := grpcServer.GetStaticFile(&staticProto.Static{Uri: staticURI}, stream)

	assert.NoError(t, err)
	stream.AssertExpectations(t)
	mockUC.AssertExpectations(t)
}

func TestGetStaticFile_Error(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)
	readSeeker := bytes.NewReader([]byte("file content"))

	staticURI := "http://example.com/static/file"
	mockUC.On("GetStaticFile", staticURI).Return(readSeeker, errors.New("file not found"))

	stream := &mockStream{}

	err := grpcServer.GetStaticFile(&staticProto.Static{Uri: staticURI}, stream)

	assert.Error(t, err)
	stream.AssertExpectations(t)
	mockUC.AssertExpectations(t)
}

func TestPing(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	result, err := grpcServer.Ping(context.Background(), &staticProto.Nothing{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

type mockStream struct {
	mock.Mock
}

func (m *mockStream) Send(resp *staticProto.StaticUpload) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *mockStream) Recv() (*staticProto.StaticUpload, error) {
	args := m.Called()
	return args.Get(0).(*staticProto.StaticUpload), args.Error(1)
}

func (m *mockStream) RecvMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *mockStream) SendAndClose(resp *staticProto.Static) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockStream) SendHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

func (m *mockStream) SendMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *mockStream) SetHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

func (m *mockStream) SetTrailer(md metadata.MD) {
	m.Called(md)
}

func (m *mockStream) CloseSend() error {
	args := m.Called()
	return args.Error(0)
}
