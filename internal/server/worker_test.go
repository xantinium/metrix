package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/server"
)

func TestWorker(t *testing.T) {
	logger.Init(true)
	defer logger.Destroy()

	incrementer := new(incrementer)
	storeInterval := 2 * time.Second
	workerLifeTime := 5 * time.Second
	expectedIncrementsNum := int(workerLifeTime.Seconds()) / int(storeInterval.Seconds())

	worker := server.NewMetrixServerWorker(storeInterval, incrementer)

	worker.Run()
	time.Sleep(workerLifeTime)
	worker.Stop()

	require.Equal(t, expectedIncrementsNum, incrementer.Counter)
}

type incrementer struct {
	Counter int
}

func (inc *incrementer) SaveMetrics(_ context.Context) error {
	inc.Counter++
	return nil
}
