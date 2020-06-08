package adapters_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZipSnapshotsArchiver_Archive(t *testing.T) {
	a := adapters.NewZipSnapshotArchiver()

	now := time.Unix(100, 0)
	then := time.Unix(200, 0)
	snaps := []application.ArchivableSnapshot{
		mockFileSnapshot{now},
		mockFileSnapshot{then},
	}

	path, err := a.Archive(snaps)
	require.NoError(t, err)
	assert.Contains(t, path, "100_200.zip")
	defer os.RemoveAll(path)
}

type mockFileSnapshot struct {
	takenAt time.Time
}

func (m mockFileSnapshot) TakenAt() time.Time {
	return m.takenAt
}

func (m mockFileSnapshot) Copy(_ io.Writer) error {
	return nil
}
