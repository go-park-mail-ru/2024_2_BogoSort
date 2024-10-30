package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

const (
	getStaticQuery = `
		SELECT path, name
		FROM static
		WHERE id = $1
	`

	uploadStaticQuery = `
		INSERT INTO static (path, name)
		VALUES ($1, $2)
		RETURNING id
	`
)

type StaticDB struct {
	DB        *pgxpool.Pool
	logger    *zap.Logger
	basicPath string
	maxSize   int
	ctx       context.Context
	timeout   time.Duration
}

func NewStaticRepository(ctx context.Context, timeout time.Duration, dbpool *pgxpool.Pool, basicPath string, maxSize int, logger *zap.Logger) (repository.StaticRepository, error) {
	if err := dbpool.Ping(ctx); err != nil {
		return nil, err
	}
	return &StaticDB{
		DB:        dbpool,
		basicPath: basicPath,
		maxSize:   maxSize,
		logger:    logger,
		ctx:       ctx,
		timeout:   timeout,
	}, nil
}

func (s StaticDB) GetStatic(staticID uuid.UUID) (string, error) {
	var path, name string

	err := s.DB.QueryRow(context.Background(), getStaticQuery, staticID).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error("postgres: static not found", zap.String("static_id", staticID.String()))
			return "", entity.PSQLWrap(repository.ErrStaticNotFound)
		}
		s.logger.Error("postgres: error getting static", zap.String("static_id", staticID.String()), zap.Error(err))
		return "", entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса GetStatic"))
	}

	s.logger.Info("postgres: static retrieved successfully", zap.String("static_id", staticID.String()), zap.String("path", path), zap.String("name", name))

	return fmt.Sprintf("%s/%s", path, name), nil
}

func (s StaticDB) UploadStatic(path, filename string, data []byte) (uuid.UUID, error) {
	if len(data) > s.maxSize {
		s.logger.Error("postgres: static too large", zap.Int("size", len(data)), zap.Int("max_size", s.maxSize))
		return uuid.UUID{}, entity.PSQLWrap(repository.ErrStaticTooLarge)
	}

	dir := filepath.Dir(fmt.Sprintf("%s/%s/", s.basicPath, path))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		s.logger.Error("error creating static directory", zap.String("path", dir), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("не удалось создать папку для хранения статики"),
		)
	}

	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename))
	if err != nil {
		s.logger.Error("error creating static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("не удалось создать файл"),
		)
	}

	if _, err = dst.Write(data); err != nil {
		s.logger.Error("error writing data to static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("не удалось записать данные в файл"),
		)
	}
	if err = dst.Close(); err != nil {
		s.logger.Error("error closing static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("не удалось закрыть файл"),
		)
	}

	var id uuid.UUID
	if err = s.DB.QueryRow(context.Background(), uploadStaticQuery, path, filename).Scan(&id); err != nil {
		s.logger.Error("error uploading static", zap.String("path", fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса UploadStatic"))
	}

	return id, nil
}
