package adapters

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/czeslavo/snappy/internal/application"
)

type ZipSnapshotsArchiver struct {
}

func NewZipSnapshotArchiver() ZipSnapshotsArchiver {
	return ZipSnapshotsArchiver{}
}

func (a ZipSnapshotsArchiver) Archive(snaps []application.ArchivableSnapshot) (path string, err error) {
	if len(snaps) < 1 {
		return "", nil
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].TakenAt().Before(snaps[j].TakenAt())
	})

	out, err := os.Create(zipPath(snaps))
	if err != nil {
		return "", fmt.Errorf("could not open file for writing an archive: %s", err)
	}

	w := zip.NewWriter(out)

	for _, snap := range snaps {
		fw, err := w.Create(snapName(snap))
		if err != nil {
			return "", fmt.Errorf("could not create snap in the archive: %s", err)
		}
		if err := snap.Copy(fw); err != nil {
			return "", fmt.Errorf("could not copy snapshot to archive: %s", err)
		}
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("could not close the archive: %s", err)
	}

	return out.Name(), nil
}

func zipPath(snaps []application.ArchivableSnapshot) string {
	first := snaps[0]
	last := snaps[len(snaps)-1]

	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf("%d_%d.zip", first.TakenAt().Unix(), last.TakenAt().Unix()),
	)
}

func snapName(snap application.ArchivableSnapshot) string {
	return fmt.Sprintf("%d.jpeg", snap.TakenAt().Unix())
}
