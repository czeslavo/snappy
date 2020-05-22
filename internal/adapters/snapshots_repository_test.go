package adapters_test

import (
	"image"
	"os"
	"testing"
	"time"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestSnapshotsFileSystemRepository(t *testing.T) {
	tempDir := os.TempDir()
	repo, err := adapters.NewSnapshotsFileSystemRepository(tempDir)
	require.NoError(t, err)

	snapshot := domain.MustNewSnapshot(image.NewGray(image.Rect(0, 0, 200, 200)), time.Now())
	err = repo.Store(snapshot)
	require.NoError(t, err)

	latest, err := repo.GetLatest()
	require.NoError(t, err)
	require.NotEmpty(t, latest)
}
