package static

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"

	"github.com/google/uuid"

	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
)

type Grpc struct {
	staticProto.UnimplementedStaticServiceServer
	staticUC usecase.StaticUseCase
}

func NewStaticGrpc(staticUC usecase.StaticUseCase) *Grpc {
	return &Grpc{staticUC: staticUC}
}

func (service *Grpc) GetStatic(_ context.Context, static *staticProto.Static) (*staticProto.Static, error) {
	staticID, err := uuid.Parse(static.GetId())
	if err != nil {
		return nil, err
	}
	uri, err := service.staticUC.GetStatic(staticID)
	if err != nil {
		return nil, err
	}
	return &staticProto.Static{Uri: uri}, nil
}

func (service *Grpc) UploadStatic(stream staticProto.StaticService_UploadStaticServer) error {
	var bytesAvatar []byte

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		bytesAvatar = append(bytesAvatar, chunk.GetChunk()...)
	}

	reader := bytes.NewReader(bytesAvatar)
	staticID, err := service.staticUC.UploadStatic(reader)
	if err != nil {
		return err
	}

	switch {
	case errors.Is(err, usecase.ErrStaticTooBigFile):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticTooBigFile"})
	case errors.Is(err, usecase.ErrStaticNotImage):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticNotImage"})
	case errors.Is(err, usecase.ErrStaticImageDimensions):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticImageDimensions"})
	case err != nil:
		return err
	default:
		return stream.SendAndClose(&staticProto.Static{Id: staticID.String()})
	}
}

func (service *Grpc) GetStaticFile(
	static *staticProto.Static,
	stream staticProto.StaticService_GetStaticFileServer,
) error {
	uri, err := url.QueryUnescape(static.GetUri())
	if err != nil {
		return err
	}
	file, err := service.staticUC.GetStaticFile(uri)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024)
	for {
		bytesCount, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&staticProto.StaticUpload{Chunk: buffer[:bytesCount]})
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *Grpc) Ping(context.Context, *staticProto.Nothing) (*staticProto.Nothing, error) {
	return &staticProto.Nothing{}, nil
}
