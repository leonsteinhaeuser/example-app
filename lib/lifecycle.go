package lib

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

type ShutdownFunc func(context.Context) error

type ProcessLifecycle interface {
	// RegisterShutdownProcess registers a process to be executed on shutdown
	// RegisterShutdownProcess must be called before Wait
	RegisterShutdownProcess(...ShutdownFunc)
	// Wait waits for the registered signals
	Wait()
	// Shutdown executes the registered processes
	// Shutown must be called after Wait
	// Shutdown blocks until all processes are executed
	Shutdown(ctx context.Context) error
}

var (
	_ ProcessLifecycle = (*processLifecycle)(nil)
)

type processLifecycle struct {
	signals   []os.Signal
	signal    chan os.Signal
	processes []ShutdownFunc
}

func NewProcessLifecycle(signals []os.Signal) *processLifecycle {
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt}
	}

	sigCh := make(chan os.Signal, 1)
	return &processLifecycle{
		signals:   signals,
		signal:    sigCh,
		processes: []ShutdownFunc{},
	}
}

// RegisterShutdownProcess registers a process to be executed on shutdown
// RegisterShutdownProcess must be called before Wait
func (l *processLifecycle) RegisterShutdownProcess(processes ...ShutdownFunc) {
	l.processes = append(l.processes, processes...)
}

// Wait waits for the registered signals
func (l *processLifecycle) Wait() {
	signal.Notify(l.signal, l.signals...)
	<-l.signal
}

func (l *processLifecycle) Shutdown(ctx context.Context) error {
	var origErr error
	for _, process := range l.processes {
		if err := process(ctx); err != nil {
			// wrap errors
			origErr = fmt.Errorf("%w: %w", origErr, err)
		}
	}
	return origErr
}
