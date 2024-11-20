package static

import (
	"context"
	static "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"strings"
	"bytes"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
	"github.com/pkg/errors"
)

const (
	bufSize = 1024
)

type StaticGrpcClient struct {
	timeout time.Duration
	staticManager static.StaticServiceClient
}

func NewStaticGrpcClient(connectAddr string, timeout time.Duration) (*StaticGrpcClient, error) {
	grpcConn, err := grpc.Dial(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	staticManager := static.NewStaticServiceClient(grpcConn)

	_, err = staticManager.Ping(context.Background(), &static.Nothing{})
	if err != nil {
		return nil, err
	}

	return &StaticGrpcClient{staticManager: staticManager, timeout: timeout}, nil
}

func (gate *StaticGrpcClient) GetStatic(staticID uuid.UUID) (string, error) {
	staticFile, err := gate.staticManager.GetStatic(context.Background(), &static.Static{Id: staticID.String()})
	if err != nil {
		if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
			return "", usecase.ErrStaticNotFound
		}
		return "", err
	}
	return staticFile.Uri, nil
}

func (gate *StaticGrpcClient) UploadStatic(reader io.ReadSeeker) (uuid.UUID, error) {
	stream, err := gate.staticManager.UploadStatic(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	zap.L().Info("Uploading static")

	buffer := make([]byte, bufSize)

	for {
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return uuid.Nil, err
		}

		err = stream.Send(&static.StaticUpload{
			Chunk: buffer[:bytesRead],
		})
		if err != nil {
			return uuid.Nil, err
		}
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		switch {
		case strings.Contains(err.Error(), usecase.ErrStaticTooBigFile.Error()):
			return uuid.Nil, usecase.ErrStaticTooBigFile
		case strings.Contains(err.Error(), usecase.ErrStaticNotImage.Error()):
			return uuid.Nil, usecase.ErrStaticNotImage
		case strings.Contains(err.Error(), usecase.ErrStaticImageDimensions.Error()):
			return uuid.Nil, usecase.ErrStaticImageDimensions
		case strings.Contains(err.Error(), "context deadline exceeded"):
			return uuid.Nil, errors.New("context deadline exceeded")
		default:
			return uuid.Nil, err
		}
	}
	zap.L().Info("Static uploaded", zap.String("id", response.Id))
	return uuid.MustParse(response.Id), nil
}

func (gate *StaticGrpcClient) GetStaticFile(staticURI string) (io.ReadSeeker, error) {
	zap.L().Info("Getting static file", zap.String("uri", staticURI))

	stream, err := gate.staticManager.GetStaticFile(context.Background(), &static.Static{Uri: staticURI})
	if err != nil {
		if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
			return nil, usecase.ErrStaticNotFound
		}
		return nil, err
	}

	var buffer []byte
	for {
		chunk, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
				return nil, usecase.ErrStaticNotFound
			}
			return nil, err
		}
		buffer = append(buffer, chunk.Chunk...)
	}

	err = stream.CloseSend()
	if err != nil {
		return nil, err
	}

	return io.ReadSeeker(bytes.NewReader(buffer)), nil
}
