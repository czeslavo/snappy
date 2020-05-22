package main

import (
	"context"
	"flag"
	"net/http"
	"sync"
	"time"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/application"
	"github.com/czeslavo/snappy/internal/ports"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logrus.New()

	cameraURL := flag.String("camera", "", "URL to camera snapshot, i.e. http://192.168.8.12/snapshot.jpeg")
	snapshotsDir := flag.String("dir", ".", "path to directory in which snapshots should be stored")
	frequency := flag.Int("freq", 30, "frequency of taking snapshots (in seconds)")

	flag.Parse()

	httpClient := &http.Client{}
	camera, err := adapters.NewJPEGCamera(httpClient, *cameraURL)
	snapshotsRepo, err := adapters.NewSnapshotsFileSystemRepository(*snapshotsDir)
	if err != nil {
		panic(err)
	}

	takeSnapshotHandler := application.NewTakeSnapshotHandler(camera, snapshotsRepo)
	latestSnapshotHandler := application.NewGetLatestSnapshotHandler(snapshotsRepo)

	var wg sync.WaitGroup

	wg.Add(2)
	ticker := ports.NewTicker(ports.TickerHandler{
		Frequency: time.Duration(*frequency) * time.Second,
		Handler:   takeSnapshotHandler.Handle,
	}, logger)
	httpServer := ports.NewHTTPServer(latestSnapshotHandler)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	defer func() {
		if err := httpServer.Close(); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := ticker.Run(ctx); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
