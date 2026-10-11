package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommandRun(t *testing.T) {
	err := rootCmd.Run(context.Background(), []string{"goilerplate"})
	assert.NoError(t, err)
}
