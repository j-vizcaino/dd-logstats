package engine

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
	return ht.AverageHitCount() > ht.AvgThreshold
}

// AverageHitCount returns the average number of hits spanning
// the history time period.
func (ht *HitTracker) AverageHitCount() uint64 {
	size := uint64(len(ht.Hits))
	// Division, rounded up
	return (ht.TotalHits + size - 1) / size
}
