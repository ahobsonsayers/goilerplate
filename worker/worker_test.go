package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorkerInvalidConfig(t *testing.T) {
	_, err := NewWorker(Config{PollingPeriod: 0})
	require.Error(t, err)
}

func TestWorkerStartStop(t *testing.T) {
	w, err := NewWorker(Config{PollingPeriod: 10 * time.Millisecond})
	require.NoError(t, err)

	done := make(chan error, 1)
	results := w.Results()
	go func() {
		done <- w.Start(context.Background())
	}()
	go func() {
		for range w.Errors() {
		}
	}()

	select {
	case result := <-results:
		assert.WithinDuration(t, time.Now(), result, time.Second)
	case <-time.After(time.Second):
		t.Fatal("worker did not emit a result")
	}

	// Already running
	err = w.Start(context.Background())
	require.Error(t, err)

	w.Stop()
	require.ErrorIs(t, <-done, context.Canceled)
	assert.False(t, w.IsRunning())

	// Channels are closed once the worker stops
	_, ok := <-results
	assert.False(t, ok)
}
