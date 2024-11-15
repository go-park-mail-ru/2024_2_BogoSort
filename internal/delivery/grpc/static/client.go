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
	"github.com/google/uuid"
)

const (
	bufSize = 1024
)

type StaticGrpcClient struct {
	staticManager static.StaticServiceClient
}

func NewStaticGrpcClient(connectAddr string) (*StaticGrpcClient, error) {
	grpcConn, err := grpc.Dial(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	staticManager := static.NewStaticServiceClient(grpcConn)

	_, err = staticManager.Ping(context.Background(), &static.Nothing{})
	if err != nil {
		return nil, err
	}

	return &StaticGrpcClient{staticManager: staticManager}, nil
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
		}
		return uuid.Nil, err
	}
	return uuid.MustParse(response.Id), nil
}
