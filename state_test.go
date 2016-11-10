package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateUpdate(t *testing.T) {
	s := State{}

	assert.Empty(t, s.Alarms)
	assert.False(t, s.AlarmIsActive)
	assert.False(t, s.IsValid())

	st := NewStats()
	update := StateUpdate{
		NewStats: st,
		AlarmState: AlarmState{
			IsActive:        false,
			AverageHitCount: 15,
		},
	}
	// No change of state → no new alarm
	s.Update(&update)
	assert.True(t, s.IsValid())
	assert.Empty(t, s.Alarms)
	assert.Equal(t, st, s.CurrentStats)

	update.AlarmState.IsActive = true
	s.Update(&update)
	// Change of state → generate an alarm
	assert.True(t, s.AlarmIsActive)
	assert.Len(t, s.Alarms, 1)
	assert.Equal(t, fmt.Sprintf(strAlarmTriggered, update.AlarmState.AverageHitCount), s.Alarms[0])

	update.AlarmState.IsActive = false
	update.AlarmState.AverageHitCount = 10
	s.Update(&update)
	// Recovery
	assert.False(t, s.AlarmIsActive)
	assert.Len(t, s.Alarms, 2)
	assert.Equal(t, fmt.Sprintf(strAlarmRecovered, update.AlarmState.AverageHitCount), s.Alarms[1])
}
