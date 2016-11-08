package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	historySize = 3
	threshold   = 5
)

func TestNewHitTracker(t *testing.T) {
	ht := NewHitTracker(historySize, threshold)

	assert.Len(t, ht.Hits, historySize)
	assert.Equal(t, 0, ht.current)
	assert.EqualValues(t, threshold, ht.AvgThreshold)
	assert.EqualValues(t, 0, ht.TotalHits)
}

func TestHitTracker(t *testing.T) {
	ht := NewHitTracker(historySize, threshold)

	ht.AddHits(2)
	ht.AddHits(8)
	ht.AddHits(5)
	assert.EqualValues(t, 15, ht.TotalHits)
	assert.False(t, ht.IsAboveThreshold())

	ht.AddHits(3)
	assert.EqualValues(t, 16, ht.TotalHits)
	assert.True(t, ht.IsAboveThreshold())

	ht.AddHits(0)
	assert.EqualValues(t, 8, ht.TotalHits)
	assert.False(t, ht.IsAboveThreshold())
}
