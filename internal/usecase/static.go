package usecase

import "github.com/google/uuid"

type StaticUseCase interface {
	// GetAvatar возвращает url аватара по id	
	GetAvatar(staticID uuid.UUID) (string, error)

	// UploadFile загружает файл и возвращает id загруженного файла
	UploadFile(data []byte) (uuid.UUID, error)

	// GetStaticURL возвращает url статики по id
	GetStaticURL(id uuid.UUID) (string, error)
}
