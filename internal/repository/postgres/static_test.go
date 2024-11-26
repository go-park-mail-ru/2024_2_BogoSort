package postgres

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	postgres2 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (pgxmock.PgxPoolIface, *postgres2.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := postgres2.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupTest(t *testing.T) (pgxmock.PgxPoolIface, *postgres2.PgxMockAdapter, string, context.Context, *StaticDB, func()) {
	tempDir := filepath.Join(os.TempDir(), "test_temp_dir")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	mockPool, adapter := setupMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &StaticDB{
		DB:        adapter,
		BasicPath: tempDir,
		MaxSize:   10 * 1024 * 1024,
		Ctx:       ctx,
		timeout:   10 * time.Second,
	}

	return mockPool, adapter, tempDir, ctx, repo, func() {
		os.RemoveAll(tempDir)
		cancel()
		mockPool.Close()
	}
}

func TestStaticDB_GetStatic(t *testing.T) {
	mockPool, _, _, _, repo, teardown := setupTest(t)
	defer teardown()

	tests := []struct {
		name          string
		staticID      uuid.UUID
		expectedPath  string
		expectedName  string
		expectedError error
	}{
		{
			name:          "Success",
			staticID:      uuid.New(),
			expectedPath:  "test/path",
			expectedName:  "test.jpg",
			expectedError: nil,
		},
		{
			name:          "NotFound",
			staticID:      uuid.New(),
			expectedPath:  "",
			expectedName:  "",
			expectedError: repository.ErrStaticNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == nil {
				rows := mockPool.NewRows([]string{"path", "name"}).AddRow(tt.expectedPath, tt.expectedName)
				mockPool.ExpectQuery("SELECT path, name FROM static WHERE id = \\$1").
					WithArgs(tt.staticID).
					WillReturnRows(rows)
			} else {
				mockPool.ExpectQuery("SELECT path, name FROM static WHERE id = \\$1").
					WithArgs(tt.staticID).
					WillReturnError(pgx.ErrNoRows)
			}

			path, err := repo.Get(tt.staticID)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "expected error in GetStatic")
			} else {
				assert.NoError(t, err, "unexpected error in GetStatic")
				expectedResult := fmt.Sprintf("%s/%s", tt.expectedPath, tt.expectedName)
				assert.Equal(t, expectedResult, path, "paths do not match")
			}

			err = mockPool.ExpectationsWereMet()
			assert.NoError(t, err, "there were unfulfilled expectations")
		})
	}
}

func TestStaticDB_UploadStatic(t *testing.T) {
	mockPool, _, tempDir, _, repo, teardown := setupTest(t)
	defer teardown()

	tests := []struct {
		name          string
		path          string
		filename      string
		data          []byte
		expectedError error
	}{
		{
			name:          "Success",
			path:          "testing/staticfiles/test/path",
			filename:      "test.jpg",
			data:          []byte("test data"),
			expectedError: nil,
		},
		{
			name:          "FileTooLarge",
			path:          "testing/staticfiles/test/path",
			filename:      "test.jpg",
			data:          bytes.Repeat([]byte("a"), 20),
			expectedError: repository.ErrStaticTooLarge,
		},
		{
			name:          "SQLError",
			path:          "testing/staticfiles/test/path",
			filename:      "test.jpg",
			data:          []byte("test data"),
			expectedError: errors.New("sql error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == nil {
				staticID := uuid.New()
				mockPool.ExpectQuery("INSERT INTO static \\(path, name\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
					WithArgs(tempDir+tt.path, tt.filename).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(staticID))
			} else if errors.Is(tt.expectedError, repository.ErrStaticTooLarge) {
				repo.MaxSize = 10
			} else {
				mockPool.ExpectQuery("INSERT INTO static \\(path, name\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
					WithArgs(tempDir+tt.path, tt.filename).
					WillReturnError(errors.New("sql error"))
			}

			_, err := repo.Upload(tt.path, tt.filename, tt.data)
			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error(), "expected error in UploadStatic")
			} else {
				assert.NoError(t, err, "unexpected error in UploadStatic")
			}

			err = mockPool.ExpectationsWereMet()
			assert.NoError(t, err, "there were unfulfilled expectations")
		})
	}
}
