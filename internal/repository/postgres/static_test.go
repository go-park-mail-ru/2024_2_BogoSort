package postgres

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	postgres2 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupMockDB(t *testing.T) (pgxmock.PgxPoolIface, *postgres2.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := postgres2.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func TestStaticDB_GetStatic_Success(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) 

	mockPool, adapter := setupMockDB(t)
	defer mockPool.Close()

	logger, _ := zap.NewDevelopment()
	repo := StaticDB{
		DB:        adapter,
		Logger:    logger,
		BasicPath: tempDir,
		MaxSize:   10 * 1024 * 1024,
		Ctx:       context.Background(),
		Timeout:   time.Second * 5,
	}

	staticID := uuid.New()
	expectedPath := "test/path"
	expectedName := "test.jpg"

	rows := mockPool.NewRows([]string{"path", "name"}).AddRow(expectedPath, expectedName)
	mockPool.ExpectQuery("SELECT path, name FROM static WHERE id = \\$1").
		WithArgs(staticID).
		WillReturnRows(rows)

	path, err := repo.GetStatic(staticID)
	assert.NoError(t, err, "unexpected error in GetStatic")

	expectedResult := fmt.Sprintf("%s/%s", expectedPath, expectedName)
	assert.Equal(t, expectedResult, path, "paths do not match")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}

func TestStaticDB_GetStatic_NotFound(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockPool, adapter := setupMockDB(t)
	defer mockPool.Close()

	logger, _ := zap.NewDevelopment()
	repo := StaticDB{
		DB:        adapter,
		Logger:    logger,
		BasicPath: tempDir,
		MaxSize:   10 * 1024 * 1024,
		Ctx:       context.Background(),
		Timeout:   time.Second * 5,
	}

	staticID := uuid.New()

	mockPool.ExpectQuery("SELECT path, name FROM static WHERE id = \\$1").
		WithArgs(staticID).
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.GetStatic(staticID)
	assert.ErrorIs(t, err, repository.ErrStaticNotFound, "expected ErrStaticNotFound in GetStatic")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}

func TestStaticDB_UploadStatic_Success(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockPool, adapter := setupMockDB(t)
	defer mockPool.Close()

	logger, _ := zap.NewDevelopment()
	repo := StaticDB{
		DB:        adapter,
		Logger:    logger,
		BasicPath: tempDir,
		MaxSize:   10 * 1024 * 1024,
		Ctx:       context.Background(),
		Timeout:   time.Second * 5,
	}

	path := "testing/staticfiles/test/path"
	filename := "test.jpg"
	data := []byte("test data")
	staticID := uuid.New()

	mockPool.ExpectQuery("INSERT INTO static \\(path, name\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(tempDir + path, filename).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(staticID))

	_, err = repo.UploadStatic(path, filename, data)
	assert.NoError(t, err, "unexpected error in UploadStatic")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}

func TestStaticDB_UploadStatic_FileTooLarge(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockPool, adapter := setupMockDB(t)
	defer mockPool.Close()

	logger, _ := zap.NewDevelopment()
	repo := StaticDB{
		DB:        adapter,
		Logger:    logger,
		BasicPath: tempDir,
		MaxSize:   10,
		Ctx:       context.Background(),
		Timeout:   time.Second * 5,
	}

	path := "testing/staticfiles/test/path"
	filename := "test.jpg"
	data := bytes.Repeat([]byte("a"), 20)

	_, err = repo.UploadStatic(path, filename, data)
	assert.ErrorIs(t, err, repository.ErrStaticTooLarge, "expected ErrStaticTooLarge in UploadStatic")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}

func TestStaticDB_UploadStatic_SQL_Error(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockPool, adapter := setupMockDB(t)
	defer mockPool.Close()

	logger, _ := zap.NewDevelopment()
	repo := StaticDB{
		DB:        adapter,
		Logger:    logger,
		BasicPath: tempDir,
		MaxSize:   10 * 1024 * 1024,
		Ctx:       context.Background(),
		Timeout:   time.Second * 5,
	}

	path := "testing/staticfiles/test/path"
	filename := "test.jpg"
	data := []byte("test data")

	mockPool.ExpectQuery("INSERT INTO static \\(path, name\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(tempDir + path, filename).
		WillReturnError(fmt.Errorf("sql error"))

	_, err = repo.UploadStatic(path, filename, data)
	assert.Error(t, err, "expected error in UploadStatic")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}