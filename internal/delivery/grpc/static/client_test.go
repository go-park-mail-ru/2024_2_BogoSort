package static

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClientStaticUseCase struct {
	mock.Mock
}

func (m *mockClientStaticUseCase) GetStatic(id uuid.UUID) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *mockClientStaticUseCase) UploadStatic(reader io.Reader) (uuid.UUID, error) {
	args := m.Called(reader)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockClientStaticUseCase) GetStaticFile(uri string) (io.Reader, error) {
	args := m.Called(uri)
	return args.Get(0).(io.Reader), args.Error(1)
}

func (m *mockClientStaticUseCase) GetAvatar(id uuid.UUID) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestClientGetStatic_Success(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
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

func TestClientGetStatic_InvalidUUID(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	static := &staticProto.Static{Id: "invalid-uuid"}

	result, err := grpcServer.GetStatic(context.Background(), static)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClientUploadStatic_Success(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	stream := &mockStream{}
	stream.On("Recv").Return(&staticProto.StaticUpload{Chunk: []byte("chunk1")}, nil).Once()
	stream.On("Recv").Return(&staticProto.StaticUpload{Chunk: []byte("chunk2")}, nil).Once()
	stream.On("Recv").Return(io.EOF)

	staticID := uuid.New()
	mockUC.On("UploadStatic", mock.Anything).Return(staticID, nil)

	err := grpcServer.UploadStatic(stream)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestClientUploadStatic_ErrorReceivingChunks(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	stream := &mockStream{}
	stream.On("Recv").Return(&staticProto.StaticUpload{Chunk: []byte("chunk1")}, nil).Once()
	stream.On("Recv").Return(errors.New("stream error"))

	err := grpcServer.UploadStatic(stream)

	assert.Error(t, err)
}

func TestClientGetStaticFile_Success(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	static := &staticProto.Static{Uri: "http://example.com/static/file"}
	stream := &mockStream{}
	stream.On("Send", &staticProto.StaticUpload{Chunk: []byte("file content")}).Return(nil)

	mockUC.On("GetStaticFile", "http://example.com/static/file").Return(bytes.NewReader([]byte("file content")), nil)

	err := grpcServer.GetStaticFile(static, stream)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestClientGetStaticFile_ErrorGettingFile(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	static := &staticProto.Static{Uri: "http://example.com/static/file"}
	stream := &mockStream{}

	mockUC.On("GetStaticFile", "http://example.com/static/file").Return(nil, errors.New("file not found"))

	err := grpcServer.GetStaticFile(static, stream)

	assert.Error(t, err)
}

func TestClientPing(t *testing.T) {
	mockUC := new(mockClientStaticUseCase)
	grpcServer := NewStaticGrpc(mockUC)

	result, err := grpcServer.Ping(context.Background(), &staticProto.Nothing{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

type mockClientStream struct {
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