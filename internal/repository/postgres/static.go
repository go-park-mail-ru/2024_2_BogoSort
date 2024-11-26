package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	DB        DBExecutor
	BasicPath string
	MaxSize   int
	Ctx       context.Context
	timeout   time.Duration
}

func (s StaticDB) GetMaxSize() int {
	return s.MaxSize
}

func NewStaticRepository(ctx context.Context, dbpool *pgxpool.Pool, basicPath string, maxSize int, logger *zap.Logger, timeout time.Duration) (repository.StaticRepository, error) {
	if err := dbpool.Ping(ctx); err != nil {
		return nil, err
	}
	return &StaticDB{
		DB:        dbpool,
		BasicPath: basicPath,
		MaxSize:   maxSize,
		Ctx:       ctx,
		timeout:   timeout,
	}, nil
}

func (s StaticDB) Get(staticID uuid.UUID) (string, error) {
	var path, name string

	ctx, cancel := context.WithTimeout(s.Ctx, s.timeout)
	defer cancel()
	logger := middleware.GetLogger(s.Ctx)
	logger.Info("getting static by id from db", zap.String("static_id", staticID.String()))

	err := s.DB.QueryRow(ctx, getStaticQuery, staticID).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("postgres: static not found", zap.String("static_id", staticID.String()))
			return "", entity.PSQLWrap(repository.ErrStaticNotFound)
		}
		logger.Error("postgres: error getting static", zap.String("static_id", staticID.String()), zap.Error(err))
		return "", entity.PSQLWrap(err, errors.New("error executing SQL query GetStatic"))
	}

	return fmt.Sprintf("%s/%s", path, name), nil
}

func (s StaticDB) Upload(path, filename string, data []byte) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(s.Ctx, s.timeout)
	defer cancel()
	logger := middleware.GetLogger(s.Ctx)
	logger.Info("uploading static to db", zap.String("path", path), zap.String("filename", filename))

	if len(data) > s.MaxSize {
		logger.Error("postgres: static too large", zap.Int("size", len(data)), zap.Int("max_size", s.MaxSize))
		return uuid.UUID{}, entity.PSQLWrap(repository.ErrStaticTooLarge)
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dir := filepath.Dir(fmt.Sprintf("%s/%s/", s.BasicPath, path))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Error("error creating static directory", zap.String("path", dir), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("failed to create a directory for storing static files"),
		)
	}

	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename))
	if err != nil {
		logger.Error("error creating static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("failed to create file"),
		)
	}

	if _, err = dst.Write(data); err != nil {
		logger.Error("error writing data to static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("failed to write data to file"),
		)
	}
	if err = dst.Close(); err != nil {
		logger.Error("error closing static file", zap.String("path", fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(
			err,
			errors.New("failed to close file"),
		)
	}

	logger.Info("static uploaded", zap.String("path", fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename)))
	var id uuid.UUID
	if err = s.DB.QueryRow(ctx, uploadStaticQuery, s.BasicPath+path, filename).Scan(&id); err != nil {
		logger.Error("error uploading static", zap.String("path", fmt.Sprintf("%s/%s/%s", s.BasicPath, path, filename)), zap.Error(err))
		return uuid.UUID{}, entity.PSQLWrap(err, errors.New("error executing SQL query UploadStatic"))
	}

	s.Logger.Info("Static file uploaded to DB", zap.String("id", id.String()))
	return id, nil
}
