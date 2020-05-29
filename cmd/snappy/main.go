package main

import (
	"context"
	"fmt"

	"github.com/czeslavo/snappy/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := service.BuildService()
	if err != nil {
		panic(fmt.Errorf("failed to build a service: %s", err))
	}

	s.Run(ctx)
}
