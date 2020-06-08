package domain

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"time"

	"github.com/h2non/filetype"
)

type FileSnapshot struct {
	path    string
	takenAt time.Time
}

func NewFileSnapshot(path string, takenAt time.Time) (FileSnapshot, error) {
	if path == "" {
		return FileSnapshot{}, errors.New("empty path")
	}
	if takenAt.IsZero() {
		return FileSnapshot{}, errors.New("empty taken at")
	}
	if err := isImage(path); err != nil {
		return FileSnapshot{}, fmt.Errorf("not an image: %s", err)
	}

	return FileSnapshot{path, takenAt}, nil
}

func isImage(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open: %s", err)
	}

	head := make([]byte, 261)
	if _, err := file.Read(head); err != nil {
		return fmt.Errorf("could not read: %s", err)
	}

	if !filetype.IsImage(head) {
		return fmt.Errorf("not an image: %s", path)
	}

	return nil
}

func (s FileSnapshot) Load() (LoadedSnapshot, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return LoadedSnapshot{}, fmt.Errorf("could not open snapshot file: %s", err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return LoadedSnapshot{}, fmt.Errorf("could not decode image: %s", err)
	}

	return NewLoadedSnapshot(img, s.takenAt)
}

func (s FileSnapshot) TakenAt() time.Time {
	return s.takenAt
}

func (s FileSnapshot) Path() string {
	return s.path
}

func (s FileSnapshot) Copy(w io.Writer) error {
	f, err := os.Open(s.path)
	if err != nil {
		return fmt.Errorf("could not open snapshot file: %s", err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("could not copy snapshot: %s", err)
	}

	return nil
}

type LoadedSnapshot struct {
	img     image.Image
	takenAt time.Time
}

func NewLoadedSnapshot(img image.Image, takenAt time.Time) (LoadedSnapshot, error) {
	if img == nil {
		return LoadedSnapshot{}, errors.New("empty image")
	}

	return LoadedSnapshot{
		img:     img,
		takenAt: takenAt,
	}, nil
}

func MustNewSnapshot(img image.Image, takenAt time.Time) LoadedSnapshot {
	s, err := NewLoadedSnapshot(img, takenAt)
	if err != nil {
		panic(err)
	}
	return s
}

func (s LoadedSnapshot) Image() image.Image {
	return s.img
}

func (s LoadedSnapshot) TakenAt() time.Time {
	return s.takenAt
}
