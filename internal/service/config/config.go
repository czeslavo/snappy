package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CameraURL string
type SnapshotsDirectory string
type SnapshotsFrequency time.Duration
type HTTPPort int

type Config struct {
	CameraURL          CameraURL
	SnapshotsDirectory SnapshotsDirectory
	SnapshotsFrequency SnapshotsFrequency
	HTTPPort           HTTPPort
}

func ReadConfigFromEnv() (Config, error) {
	cameraURL := requireEnv("CAMERA_URL")
	snapshotsDir := getEnv("SNAPSHOTS_DIR", ".")
	snapshotsFrequency := getInt("SNAPSHOTS_FREQUENCY", 30)
	httpPort := getInt("HTTP_PORT", 8080)

	return Config{
		CameraURL:          CameraURL(cameraURL),
		SnapshotsDirectory: SnapshotsDirectory(snapshotsDir),
		SnapshotsFrequency: SnapshotsFrequency(time.Second * time.Duration(snapshotsFrequency)),
		HTTPPort:           HTTPPort(httpPort),
	}, nil
}

func requireEnv(key string) string {
	if val := os.Getenv(key); val == "" {
		panic(fmt.Errorf("missing %s environment variable", key))
	} else {
		return val
	}
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if frequency := os.Getenv(key); frequency != "" {
		i, err := strconv.Atoi(frequency)
		if err != nil {
			panic(fmt.Errorf("expected %s to be int: %s", key, err))
		}

		return i
	}

	return defaultValue
}
