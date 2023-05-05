package worker

import (
	"context"
)

// Worker is a function that starts a number of workers in parallel and waits for them to finish
// The workerFunc is called in a loop for each worker until the context is cancelled
func Worker(ctx context.Context, workerCount int, workerFunc func(ctx context.Context, workerID int)) {
	for i := 0; i < workerCount; i++ {
		go workerFunc(ctx, i)
	}
}
