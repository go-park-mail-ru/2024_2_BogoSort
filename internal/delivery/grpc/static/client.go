package static

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"strings"
	"time"

	static "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	bufSize = 1024
)

type StaticGrpcClient struct {
	timeout       time.Duration
	staticManager static.StaticServiceClient
}

func NewStaticGrpcClient(connectAddr string, timeout time.Duration) (*StaticGrpcClient, error) {
	//nolint:staticcheck // Suppressing deprecation warning for grpc.Dial
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := gate.staticManager.UploadStatic(ctx)
	if err != nil {
		zap.L().Error("Ошибка при инициализации потока UploadStatic", zap.Error(err))
		return uuid.Nil, err
	}

	zap.L().Info("Начало загрузки статического файла")

	imgData, err := io.ReadAll(reader)
	if err != nil {
		zap.L().Error("Ошибка чтения данных изображения", zap.Error(err))
		return uuid.Nil, err
	}

	cleanImgData, err := removeMetadata(imgData)
	if err != nil {
		zap.L().Error("Ошибка удаления метаданных", zap.Error(err))
		return uuid.Nil, err
	}

	buffer := bytes.NewBuffer(cleanImgData)

	chunk := make([]byte, bufSize)
	for {
		bytesRead, err := buffer.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			zap.L().Error("Ошибка чтения чанка", zap.Error(err))
			return uuid.Nil, err
		}

		err = stream.Send(&static.StaticUpload{
			Chunk: chunk[:bytesRead],
		})
		if err != nil {
			zap.L().Error("Ошибка отправки чанка", zap.Error(err))
			return uuid.Nil, err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		zap.L().Error("Ошибка при закрытии и получении ответа", zap.Error(err))
		// Обработка ошибок, как ранее
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

	zap.L().Info("Статический файл успешно загружен", zap.String("id", response.Id))
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

func removeMetadata(imgData []byte) ([]byte, error) {
	var (
		img image.Image
		err error
	)

	zap.L().Info("Начало удаления метаданных из изображения")

	// Попытка декодировать изображение стандартными методами
	img, _, err = image.Decode(bytes.NewReader(imgData))
	if err != nil {
		zap.L().Error("Не удалось декодировать изображение стандартными методами", zap.Error(err))
		return nil, err
	}

	var buf bytes.Buffer

	// Кодирование изображения обратно в JPEG без метаданных
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		zap.L().Error("Ошибка кодирования изображения в JPEG", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Метаданные успешно удалены")
	return buf.Bytes(), nil
}
