// Package worker provides an example background worker, showing the
// boilerplate for a goroutine based worker with proper synchronisation.
//
// Replace Worker.Run() with your own code. If your project does not need a
// worker, delete this package and the worker command in cmd/goilerplate.
package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

func NewWorker(config Config) (*Worker, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	return &Worker{
		config: config,
		// Result Channels
		resultChan: make(chan time.Time),
		errorChan:  make(chan error),
	}, nil
}

type Config struct {
	PollingPeriod time.Duration
}

func (c *Config) Validate() error {
	if c.PollingPeriod <= 0 {
		return errors.New("polling period must be a valid positive duration")
	}

	return nil
}

type Worker struct {
	config Config

	// Result Channels
	resultChan chan time.Time
	errorChan  chan error

	// Synchronisation
	running      atomic.Bool
	runningWg    sync.WaitGroup
	shutdownOnce sync.Once

	// Cancellation
	cancel context.CancelFunc
}

func (w *Worker) IsRunning() bool           { return w.running.Load() }
func (w *Worker) WaitUntilStopped()         { w.runningWg.Wait() }
func (w *Worker) Results() <-chan time.Time { return w.resultChan }
func (w *Worker) Errors() <-chan error      { return w.errorChan }

// Start the worker - blocks until stopped
func (w *Worker) Start(ctx context.Context) error {
	swapped := w.running.CompareAndSwap(false, true)
	if !swapped {
		return errors.New("worker already running")
	}
	w.runningWg.Add(1)
	defer w.cleanup()

	// Create ticker
	ticker := time.NewTicker(w.config.PollingPeriod)
	defer ticker.Stop()

	// Create cancellable context
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	// Start worker loop
	for {
		select {
		case <-ticker.C:
			w.Process()

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Process is the worker process
func (w *Worker) Process() {
	currentTime := time.Now()

	w.resultChan <- currentTime

	// Example error emission
	if currentTime.Second()%5 == 0 {
		w.errorChan <- errors.New("example worker error")
	}
}

func (w *Worker) cleanup() {
	w.shutdownOnce.Do(func() {
		// Close channels
		close(w.resultChan)
		close(w.errorChan)

		// Signal we are no longer running
		w.running.Store(false)
		w.runningWg.Done()
	})
}

// Stop the worker - blocks until stopped
func (w *Worker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.WaitUntilStopped()
}
