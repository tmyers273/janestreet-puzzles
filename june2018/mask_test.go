package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItWorks(t *testing.T) {
	key := 15697447
	mask := uint64(391512286756864)

	final := keyToMask(uint64(key))

	assert.Equal(t, mask, final)
}
