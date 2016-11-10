package ui

import (
	"dd-logstats/engine"
	"fmt"
)

type State struct {
	CurrentStats  *engine.Stats
	Alarms        []string
	AlarmIsActive bool
}

type AlarmState struct {
	IsActive        bool
	AverageHitCount uint64
}

type StateUpdate struct {
	NewStats   *engine.Stats
	AlarmState AlarmState
}

const (
	strAlarmTriggered = "High traffic generated an alarm (avg hits=%d)"
	strAlarmRecovered = "Traffic returned to normal (avg hits=%d)"
)

func (s *State) Update(u *StateUpdate) {
	s.CurrentStats = u.NewStats

	if s.AlarmIsActive == u.AlarmState.IsActive {
		return
	}
	var msg string
	if u.AlarmState.IsActive {
		msg = fmt.Sprintf(strAlarmTriggered, u.AlarmState.AverageHitCount)
	} else {
		msg = fmt.Sprintf(strAlarmRecovered, u.AlarmState.AverageHitCount)
	}
	s.Alarms = append(s.Alarms, msg)
	s.AlarmIsActive = u.AlarmState.IsActive
}

func (s *State) IsValid() bool {
	return s.CurrentStats != nil
}
