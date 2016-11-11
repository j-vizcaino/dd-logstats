package ui

import (
	"dd-logstats/engine"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStateUpdate(t *testing.T) {
	s := State{}

	assert.Empty(t, s.Alarms)
	assert.False(t, s.AlarmIsActive)
	assert.False(t, s.IsValid())

	st := engine.NewStats()
	st.Finalize()
	// No change of state → no new alarm
	s.Update(st, false, 10)
	assert.True(t, s.IsValid())
	assert.Empty(t, s.Alarms)
	assert.Equal(t, st, s.CurrentStats)

	s.Update(st, true, 15)
	// Change of state → generate an alarm
	assert.True(t, s.AlarmIsActive)
	assert.Len(t, s.Alarms, 1)
	assert.Equal(t, Alarm{Timestamp: st.DateEnd, Active: true, AverageHits: 15}, s.Alarms[0])

	s.Update(st, true, 18)
	assert.True(t, s.AlarmIsActive)
	// Did not generate a new alarm
	assert.Len(t, s.Alarms, 1)

	st.DateEnd = time.Now()
	s.Update(st, false, 10)
	// Recovery
	assert.False(t, s.AlarmIsActive)
	assert.Len(t, s.Alarms, 2)
	assert.Equal(t, Alarm{Timestamp: st.DateEnd, Active: false, AverageHits: 10}, s.Alarms[1])
}
