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

func (r SnapshotsFileSystemRepository) Store(snapshot domain.LoadedSnapshot) (err error) {
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

func (r SnapshotsFileSystemRepository) GetAll() ([]domain.FileSnapshot, error) {
	files, err := filepath.Glob(r.rootDir + "/*.jpeg")
	if err != nil {
		return nil, fmt.Errorf("failed to glob files in root directory: %s", err)
	}

	var snaps []domain.FileSnapshot
	for _, f := range files {
		filename := filepath.Base(f)
		takenTime, err := r.extractTakenTimeFromSnapshotFilename(filename)
		if err != nil {
			return nil, fmt.Errorf("could not extract taken time: %s", err)
		}
		snap, err := domain.NewFileSnapshot(f, takenTime)
		if err != nil {
			return nil, fmt.Errorf("could not create file snapshot: %s", err)
		}

		snaps = append(snaps, snap)
	}

	return snaps, nil
}

func (r SnapshotsFileSystemRepository) GetLatest() (domain.LoadedSnapshot, error) {
	files, err := filepath.Glob(r.rootDir + "/*.jpeg")
	if err != nil {
		return domain.LoadedSnapshot{}, fmt.Errorf("failed to glob files in root directory: %s", err)
	}
	if len(files) < 1 {
		return domain.LoadedSnapshot{}, errors.New("no latest snapshot available")
	}

	sort.Strings(files)

	latestPath := files[len(files)-1]
	latestFilename := filepath.Base(latestPath)
	takenTime, err := r.extractTakenTimeFromSnapshotFilename(latestFilename)
	if err != nil {
		return domain.LoadedSnapshot{}, fmt.Errorf("could not extract taken time from filename: %s", err)
	}

	snap, err := domain.NewFileSnapshot(latestPath, takenTime)
	if err != nil {
		return domain.LoadedSnapshot{}, fmt.Errorf("could not create file snapshot: %s", err)
	}

	return snap.Load()
}

func (r SnapshotsFileSystemRepository) snapshotPath(s domain.LoadedSnapshot) string {
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
	snapshots []domain.LoadedSnapshot
}

func (r *SnapshotsInMemoryRepository) Store(snapshot domain.LoadedSnapshot) error {
	r.snapshots = append(r.snapshots, snapshot)
	return nil
}

func (r *SnapshotsInMemoryRepository) GetLatest() (domain.LoadedSnapshot, error) {
	if len(r.snapshots) < 1 {
		return domain.LoadedSnapshot{}, errors.New("no latest snapshot available")
	}

	return r.snapshots[len(r.snapshots)-1], nil
}
