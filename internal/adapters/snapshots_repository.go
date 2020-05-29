package adapters

import (
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/czeslavo/snappy/internal/service/config"

	"github.com/czeslavo/snappy/internal/domain"
)

type SnapshotsFileSystemRepository struct {
	rootDir string
}

func NewSnapshotsFileSystemRepository(snapshotsDir config.SnapshotsDirectory) (SnapshotsFileSystemRepository, error) {
	rootDir := string(snapshotsDir)
	info, err := os.Stat(rootDir)
	if err != nil {
		return SnapshotsFileSystemRepository{}, fmt.Errorf("could not stat root dir: %s", err)
	}
	if !info.IsDir() {
		return SnapshotsFileSystemRepository{}, fmt.Errorf("'%s' is not directory", rootDir)
	}

	return SnapshotsFileSystemRepository{rootDir}, nil
}

func (r SnapshotsFileSystemRepository) Store(snapshot domain.Snapshot) (err error) {
	path := r.snapshotPath(snapshot)
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %s", err)
	}
	defer func() {
		err = out.Close()
	}()

	if err := jpeg.Encode(out, snapshot.Image(), nil); err != nil {
		return fmt.Errorf("could not encode snapshot to jpeg: %s", err)
	}

	return nil
}

func (r SnapshotsFileSystemRepository) GetLatest() (domain.Snapshot, error) {
	files, err := filepath.Glob(r.rootDir + "/*.jpeg")
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("failed to glob files in root directory: %s", err)
	}
	if len(files) < 1 {
		return domain.Snapshot{}, errors.New("no latest snapshot available")
	}

	sort.Strings(files)

	latestPath := files[len(files)-1]
	latest, err := os.Open(latestPath)
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("could not open latest snapshot file: %s", err)
	}
	defer latest.Close()

	img, err := jpeg.Decode(latest)
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("could not decode image: %s", err)
	}

	latestFilename := filepath.Base(latestPath)
	takenTime, err := r.extractTakenTimeFromSnapshotFilename(latestFilename)
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("could not extract taken time from filename: %s", err)
	}

	return domain.NewSnapshot(img, takenTime)
}

func (r SnapshotsFileSystemRepository) snapshotPath(s domain.Snapshot) string {
	filename := strconv.Itoa(int(s.TakenAt().Unix())) + ".jpeg"
	return filepath.Join(r.rootDir, filename)
}

func (r SnapshotsFileSystemRepository) extractTakenTimeFromSnapshotFilename(filename string) (time.Time, error) {
	withoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	unixSec, err := strconv.ParseInt(withoutExt, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not convert filename to unix seconds: %s", err)
	}

	return time.Unix(unixSec, 0), nil
}

type SnapshotsInMemoryRepository struct {
	snapshots []domain.Snapshot
}

func (r *SnapshotsInMemoryRepository) Store(snapshot domain.Snapshot) error {
	r.snapshots = append(r.snapshots, snapshot)
	return nil
}

func (r *SnapshotsInMemoryRepository) GetLatest() (domain.Snapshot, error) {
	if len(r.snapshots) < 1 {
		return domain.Snapshot{}, errors.New("no latest snapshot available")
	}

	return r.snapshots[len(r.snapshots)-1], nil
}
