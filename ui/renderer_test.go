package ui

import (
	"testing"
	"time"

	"dd-logstats/engine"

	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, s *State) string {
	r, err := NewRenderer("assets", 10*time.Second, 100, time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Result())

	err = r.Render(s)
	assert.NoError(t, err)

	return r.Result()
}

const dateFormat = "2006-01-02 15:04:05"

func TestRendererInitial(t *testing.T) {
	s := render(t, &State{})

	assert.Regexp(t, `\bno stats yet\b`, s)
	assert.Regexp(t, `\bNo traffic alarm so far\b`, s)
	assert.Regexp(t, `Traffic is normal`, s)
}

func TestRendererAlarms(t *testing.T) {

	ts1 := time.Now().UTC()
	ts2 := ts1.Add(1 * time.Second)
	s := render(t, &State{
		AlarmIsActive: true,
		Alarms: []Alarm{
			Alarm{
				Timestamp:   ts1,
				Active:      false,
				AverageHits: 10,
			},
			Alarm{
				Timestamp:   ts2,
				Active:      true,
				AverageHits: 15,
			},
		},
	})

	assert.Regexp(t, `\bHigh traffic\b`, s)
	assert.Regexp(t, ts1.Format(dateFormat), s)
	assert.Regexp(t, ts2.Format(dateFormat), s)
}

func TestRendererHits(t *testing.T) {
	logEntries := []*engine.LogEntry{
		&engine.LogEntry{
			ClientIP:  "1.1.1.1",
			URL:       "/foo/1",
			HTTP:      engine.HTTPInfo{ReturnCode: 200},
			SizeBytes: 123,
		},
		&engine.LogEntry{
			ClientIP:  "1.1.1.1",
			URL:       "/foo/2",
			HTTP:      engine.HTTPInfo{ReturnCode: 404},
			SizeBytes: 321,
		},
		&engine.LogEntry{
			ClientIP:  "1.1.1.2",
			URL:       "/foo/1",
			HTTP:      engine.HTTPInfo{ReturnCode: 201},
			SizeBytes: 987,
		},
		&engine.LogEntry{
			ClientIP:  "1.1.2.1",
			URL:       "/bar/1",
			HTTP:      engine.HTTPInfo{ReturnCode: 301},
			SizeBytes: 12,
		},
		&engine.LogEntry{
			ClientIP:  "1.1.2.1",
			URL:       "/bar/2",
			HTTP:      engine.HTTPInfo{ReturnCode: 404},
			SizeBytes: 4096,
		},
	}
	stats := engine.NewStats()
	for _, log := range logEntries {
		stats.AddLogEntry(log)
	}

	s := render(t, &State{CurrentStats: stats})

	assert.Regexp(t, `\bfoo\b`, s)
	assert.Regexp(t, `\bbar\b`, s)
	assert.Regexp(t, `\b1\.1\.1\.1\b`, s)
	assert.Regexp(t, `\b1\.1\.1\.2\b`, s)
	assert.Regexp(t, `\b1\.1\.2\.1\b`, s)
	assert.Regexp(t, `\b2xx\b`, s)
	assert.Regexp(t, `\b3xx\b`, s)
	assert.Regexp(t, `\b4xx\b`, s)
}
