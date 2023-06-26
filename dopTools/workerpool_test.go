package dopTools

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkerPool(t *testing.T) {
	workerCount := 2
	taskCount := 5

	// Define a worker function that sleeps for 100ms and returns the input value
	workerFun := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
		case <-time.After(100 * time.Millisecond):
		}
		return nil
	}

	wp := NewWorkerPool(context.Background(), workerCount, 0)
	for i := 0; i < taskCount; i++ {
		wp.Submit(workerFun)
	}
	require.NoError(t, wp.FinishAndWait(), "worker pool failed with error")
	require.True(t, wp.GetDuration() > 285*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 315*time.Millisecond, "worker pool finished too slow")

	wp = NewWorkerPool(context.Background(), 1, 0)
	for i := 0; i < taskCount; i++ {
		wp.Submit(workerFun)
	}
	require.NoError(t, wp.FinishAndWait(), "worker pool failed with error")
	require.True(t, wp.GetDuration() > 485*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 515*time.Millisecond, "worker pool finished too slow")

	wp = NewWorkerPool(context.Background(), workerCount, 0)
	go func() {
		defer wp.Finish()
		for i := 0; i < taskCount; i++ {
			wp.Submit(workerFun)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	wp.Cancel(nil)
	require.NoError(t, wp.Wait(), "worker pool failed with error")
	fmt.Println(wp.GetDuration())
	require.True(t, wp.GetDuration() > 150*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 160*time.Millisecond, "worker pool finished too slow")
}
