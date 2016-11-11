package ui

import (
	"dd-logstats/engine"
	"time"
)

type Alarm struct {
	Timestamp   time.Time
	AverageHits uint64
	Active      bool
}
type State struct {
	CurrentStats  *engine.Stats
	Alarms        []Alarm
	AlarmIsActive bool
}

func (s *State) Update(stats *engine.Stats, alarmIsActive bool, averageHitCount uint64) {
	s.CurrentStats = stats

	if s.AlarmIsActive == alarmIsActive {
		return
	}
	s.Alarms = append(s.Alarms, Alarm{
		Timestamp:   s.CurrentStats.DateEnd.UTC(),
		Active:      alarmIsActive,
		AverageHits: averageHitCount,
	})
	s.AlarmIsActive = alarmIsActive
}

func (s *State) IsValid() bool {
	return s.CurrentStats != nil
}
