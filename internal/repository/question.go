package repository

import "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"

type QuestionRepository interface {
	Create(question *entity.Question) error
	GetByPageType(pageType entity.PageType) ([]entity.Question, error)
	GetAll() ([]entity.Question, error)
}
