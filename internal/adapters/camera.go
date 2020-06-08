package adapters

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/czeslavo/snappy/internal/service/config"

	"github.com/czeslavo/snappy/internal/domain"
)

type JPEGCamera struct {
	client *http.Client
	url    string
}

func NewJPEGCamera(client *http.Client, url config.CameraURL) (JPEGCamera, error) {
	if client == nil {
		return JPEGCamera{}, errors.New("empty client")
	}
	if url == "" {
		return JPEGCamera{}, errors.New("empty url")
	}

	return JPEGCamera{
		client: client,
		url:    string(url),
	}, nil
}

func (c JPEGCamera) Get() (domain.LoadedSnapshot, error) {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return domain.LoadedSnapshot{}, fmt.Errorf("could not get snapshot: %s", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return domain.LoadedSnapshot{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return domain.LoadedSnapshot{}, fmt.Errorf("could not decode image: %s", err)
	}

	return domain.NewLoadedSnapshot(img, time.Now())
}
