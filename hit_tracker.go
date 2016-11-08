package main

// HitTracker tracks n number of hits.
type HitTracker struct {
	Hits         []uint64
	TotalHits    uint64
	AvgThreshold uint64
	current      int
}

// NewHitTracker creates a new HitTracker
// - historySize caps the total number of hits that can be stored
// - avgThreshold alarm when TotalHits/historySize is greater than this number
func NewHitTracker(historySize, avgThreshold uint64) *HitTracker {
	return &HitTracker{
		Hits:         make([]uint64, historySize),
		AvgThreshold: uint64(avgThreshold),
	}
}

// AddHits adds nbHits to the tracker
func (ht *HitTracker) AddHits(nbHits uint64) {
	ht.TotalHits -= ht.Hits[ht.current]
	ht.TotalHits += nbHits
	ht.Hits[ht.current] = nbHits
	ht.current = (ht.current + 1) % len(ht.Hits)
}

// IsAboveThreshold checks if the total number of hits are greater than
// the maximum threshold.
func (ht *HitTracker) IsAboveThreshold() bool {
	size := uint64(len(ht.Hits))
	// Division, rounded up
	avg := (ht.TotalHits + size - 1) / size
	return avg > ht.AvgThreshold
}
