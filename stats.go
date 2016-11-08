package main

import (
	"fmt"
	"sort"
)

// SectionStats holds statistics for a given URL section
type SectionStats struct {
	SectionName    string
	HitCount       uint64
	ReturnedCode   map[string]uint
	BytesPerClient map[string]uint64
}

// Helper function: transform HTTP return code into code family (1xx, 2xx, 3xx, etc...)
func codeFamily(returnCode uint) string {
	c := returnCode / 100
	return fmt.Sprintf("%dxx", c)
}

func newSectionStats(name string) *SectionStats {
	return &SectionStats{
		SectionName:    name,
		ReturnedCode:   make(map[string]uint),
		BytesPerClient: make(map[string]uint64),
	}
}

// AddLogEntry increases the current SectionStats statistics
// with the give LogEntry
func (s *SectionStats) AddLogEntry(l *LogEntry) {
	s.HitCount++
	k := codeFamily(l.HTTP.ReturnCode)
	s.ReturnedCode[k]++
	s.BytesPerClient[l.ClientIP] += l.SizeBytes
}

// Stats holds statistics grouped by URL section
type Stats struct {
	Sections  []*SectionStats
	ByName    map[string]*SectionStats
	TotalHits uint64
}

// NewStats creates a new empty Stats object
func NewStats() *Stats {
	return &Stats{
		ByName: make(map[string]*SectionStats),
	}
}

// AddLogEntry adds a new LogEntry to the stats
func (h *Stats) AddLogEntry(l *LogEntry) {
	section := l.URLSection()

	s, ok := h.ByName[section]
	if !ok {
		s = newSectionStats(section)
		h.Sections = append(h.Sections, s)
		h.ByName[section] = s
	}
	s.AddLogEntry(l)
	h.TotalHits++
}

// SectionsByHits returns a slice of SectionStats sorted by hit count.
// The section with the most hits comes first, the one with the least, last
func (h *Stats) SectionsByHits() []*SectionStats {
	moreHits := func(a, b *SectionStats) bool {
		return a.HitCount > b.HitCount
	}
	stats := h.Sections
	by(moreHits).Sort(stats)
	return stats
}

//
// Custom sort helpers
//

type by func(a, b *SectionStats) bool

func (b by) Sort(s []*SectionStats) {
	sorter := &sectionStatSorter{
		stats:   s,
		cmpFunc: b,
	}
	sort.Sort(sorter)
}

type sectionStatSorter struct {
	stats   []*SectionStats
	cmpFunc by
}

func (s *sectionStatSorter) Len() int {
	return len(s.stats)
}

func (s *sectionStatSorter) Swap(i, j int) {
	s.stats[i], s.stats[j] = s.stats[j], s.stats[i]
}

func (s *sectionStatSorter) Less(i, j int) bool {
	return s.cmpFunc(s.stats[i], s.stats[j])
}
