package agent_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xantinium/metrix/internal/agent"
	"github.com/xantinium/metrix/internal/logger"
)

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Init(true)
	defer logger.Destroy()

	counter, uploadFunc := getCounter()

	poolSize := 5
	reportInterval := 2 * time.Second
	workerLifeTime := 5 * time.Second
	expectedIncrementsNum := (int(workerLifeTime.Seconds()) / int(reportInterval.Seconds())) * poolSize

	worker := agent.NewMetrixAgentWorkerPool(agent.MetrixAgentWorkerPoolOptions{
		PoolSize:       poolSize,
		ReportInterval: reportInterval,
		UploadFunc:     uploadFunc,
	})

	worker.Run(ctx)
	time.Sleep(workerLifeTime)
	cancel()

	require.Equal(t, int32(expectedIncrementsNum), counter.Load())
}

func getCounter() (*atomic.Int32, func()) {
	counter := new(atomic.Int32)

	return counter, func() {
		counter.Add(1)
	}
}
