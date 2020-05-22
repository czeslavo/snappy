package adapters_test

import (
	"context"
	"image"
	"net/http"
	"testing"
	"time"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJPEGCamera(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go serveSampleJPEG(t, ctx)

	client := &http.Client{}
	c, err := adapters.NewJPEGCamera(client, "http://localhost/sample.jpeg")
	require.NoError(t, err)

	snap, err := c.Get()
	require.NoError(t, err)
	require.NotNil(t, snap)

	assert.Equal(t, image.Point{X: 0, Y: 0}, snap.Image().Bounds().Min)
	assert.Equal(t, image.Point{X: 1920, Y: 1080}, snap.Image().Bounds().Max)
	assert.NotEmpty(t, snap.TakenAt())
}

func serveSampleJPEG(t *testing.T, ctx context.Context) {
	s := &http.Server{
		Handler: http.FileServer(http.Dir("./")),
	}
	go func() {
		<-ctx.Done()
		require.NoError(t, s.Shutdown(ctx))
	}()
	_ = s.ListenAndServe()
}
