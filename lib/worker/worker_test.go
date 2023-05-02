package worker

import (
	"context"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	type args struct {
		ctxTimeoutSec int
		workerCount   int
		workerFunc    func(ctx context.Context, workerID int)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with timeout",
			args: args{
				ctxTimeoutSec: 10,
				workerCount:   5,
				workerFunc: func(ctx context.Context, workerID int) {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							time.Sleep(1 * time.Second)
						}
					}
				},
			},
		},
		{
			name: "without timeout",
			args: args{
				ctxTimeoutSec: 0,
				workerCount:   5,
				workerFunc: func(ctx context.Context, workerID int) {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							time.Sleep(1 * time.Second)
						}
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.ctxTimeoutSec > 0 {
				ctx2, cf := context.WithTimeout(ctx, time.Duration(tt.args.ctxTimeoutSec))
				defer cf()
				ctx = ctx2
			}
			Worker(ctx, tt.args.workerCount, tt.args.workerFunc)
		})
	}
}
