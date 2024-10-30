package usecase

import "github.com/google/uuid"

type Static interface {
	// GetAvatar возвращает url аватара по id	
	GetAvatar(staticID uuid.UUID) (string, error)

	// UploadAvatar загружает аватар и возвращает id загруженного файла
	UploadAvatar(data []byte) (uuid.UUID, error)

	// GetStaticURL возвращает url статики по id
	GetStaticURL(id uuid.UUID) (string, error)
}
