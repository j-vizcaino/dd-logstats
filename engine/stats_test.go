package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeFamily(t *testing.T) {
	assert.Equal(t, "1xx", codeFamily(105))
	assert.Equal(t, "4xx", codeFamily(404))
	// does not care about return code range
	assert.Equal(t, "9xx", codeFamily(999))
}

var fooEntries = []*LogEntry{
	&LogEntry{
		ClientIP:  "1.1.1.1",
		URL:       "/foo/1",
		HTTP:      HTTPInfo{ReturnCode: 200},
		SizeBytes: 123,
	},
	&LogEntry{
		ClientIP:  "1.1.1.1",
		URL:       "/foo/2",
		HTTP:      HTTPInfo{ReturnCode: 404},
		SizeBytes: 321,
	},
	&LogEntry{
		ClientIP:  "1.1.1.2",
		URL:       "/foo/1",
		HTTP:      HTTPInfo{ReturnCode: 201},
		SizeBytes: 987,
	},
}

func testFooSection(t *testing.T, s *SectionStats) {
	expectedBytes := map[string]uint64{
		"1.1.1.1": uint64(123 + 321),
		"1.1.1.2": uint64(987),
	}
	assert.Equal(t, "foo", s.SectionName)
	assert.EqualValues(t, expectedBytes, s.BytesPerClient)
	assert.Equal(t, map[string]uint{"2xx": 2, "4xx": 1}, s.ReturnedCode)
	assert.EqualValues(t, len(fooEntries), s.HitCount)
}

var barEntries = []*LogEntry{
	&LogEntry{
		ClientIP:  "1.1.2.1",
		URL:       "/bar/1",
		HTTP:      HTTPInfo{ReturnCode: 301},
		SizeBytes: 12,
	},
	&LogEntry{
		ClientIP:  "1.1.2.1",
		URL:       "/bar/2",
		HTTP:      HTTPInfo{ReturnCode: 404},
		SizeBytes: 4096,
	},
}

func testBarSection(t *testing.T, s *SectionStats) {
	expectedBytes := map[string]uint64{
		"1.1.2.1": uint64(12 + 4096),
	}
	assert.Equal(t, "bar", s.SectionName)
	assert.EqualValues(t, expectedBytes, s.BytesPerClient)
	assert.Equal(t, map[string]uint{"3xx": 1, "4xx": 1}, s.ReturnedCode)
	assert.EqualValues(t, len(barEntries), s.HitCount)
}

func TestSectionStatsAddLogEntry(t *testing.T) {
	s := newSectionStats("foo")

	for _, f := range fooEntries {
		s.AddLogEntry(f)
	}
	testFooSection(t, s)
}

func TestStatsAddLogEntry(t *testing.T) {
	s := NewStats()
	assert.True(t, s.DateEnd.IsZero())
	assert.False(t, s.DateStart.IsZero())

	for _, l := range barEntries {
		s.AddLogEntry(l)
	}
	assert.Len(t, s.Sections, 1)
	assert.Len(t, s.ByName, 1)
	assert.Equal(t, s.Sections[0], s.ByName["bar"])
	assert.EqualValues(t, len(barEntries), s.TotalHits)
	testBarSection(t, s.ByName["bar"])

	for _, l := range fooEntries {
		s.AddLogEntry(l)
	}
	assert.Len(t, s.Sections, 2)
	assert.Len(t, s.ByName, 2)
	assert.EqualValues(t, len(barEntries)+len(fooEntries), s.TotalHits)
	assert.NotNil(t, s.ByName["foo"])
	testFooSection(t, s.ByName["foo"])

	hits := s.SectionsByHits()
	assert.Len(t, hits, 2)
	// Sort by hits, DESC: expect [foo, bar]
	assert.Equal(t, hits[0], s.ByName["foo"])
	assert.Equal(t, hits[1], s.ByName["bar"])
}
