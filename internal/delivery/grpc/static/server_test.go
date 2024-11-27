package static

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
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

func (m *mockStaticUseCase) UploadStatic(reader io.Reader) (uuid.UUID, error) {
	args := m.Called(reader)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockStaticUseCase) GetStaticFile(uri string) (io.Reader, error) {
	args := m.Called(uri)
	return args.Get(0).(io.Reader), args.Error(1)
}

func TestGetStatic(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	staticID := uuid.New()
	expectedURI := "http://example.com/static/" + staticID.String()

	mockUC.On("GetStatic", staticID).Return(expectedURI, nil)

	static := &staticProto.Static{Id: staticID.String()}

	result, err := grpcServer.GetStatic(context.Background(), static)

	assert.NoError(t, err)
	assert.Equal(t, expectedURI, result.Uri)
	mockUC.AssertExpectations(t)
}

func TestUploadStatic(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	stream := &mockStream{}
	// Simulate receiving chunks
	stream.On("Recv").Return(&staticProto.StaticUpload{Chunk: []byte("chunk1")}, nil).Once()
	stream.On("Recv").Return(&staticProto.StaticUpload{Chunk: []byte("chunk2")}, nil).Once()
	stream.On("Recv").Return(io.EOF)

	mockUC.On("UploadStatic", mock.Anything).Return(uuid.New(), nil)

	err := grpcServer.UploadStatic(stream)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestGetStaticFile(t *testing.T) {
	mockUC := new(mockStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	static := &staticProto.Static{Uri: "http://example.com/static/file"}
	stream := &mockStream{}
	// Simulate sending chunks
	stream.On("Send", &staticProto.StaticUpload{Chunk: []byte("file content")}).Return(nil)

	mockUC.On("GetStaticFile", "http://example.com/static/file").Return(bytes.NewReader([]byte("file content")), nil)

	err := grpcServer.GetStaticFile(static, stream)

	assert.NoError(t, err)
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

func (m *mockStream) Recv() (*staticProto.StaticUpload, error) {
	args := m.Called()
	return args.Get(0).(*staticProto.StaticUpload), args.Error(1)
}