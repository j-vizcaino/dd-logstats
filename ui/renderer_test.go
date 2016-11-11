package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, s *State) string {
	r, err := NewRenderer("assets")
	assert.NoError(t, err)

	err = r.Render(s)
	assert.NoError(t, err)

	return r.Result()
}

func TestRendererAlarms(t *testing.T) {

	ts1 := time.Now()
	ts2 := ts1.Add(1 * time.Second)
	render(t, &State{
		AlarmIsActive: false,
		Alarms: []Alarm{
			Alarm{
				Timestamp:   ts1,
				Active:      true,
				AverageHits: 15,
			},
			Alarm{
				Timestamp:   ts2,
				Active:      false,
				AverageHits: 10,
			},
		},
	})

}

func TestRendererInitial(t *testing.T) {
	s := render(t, &State{})

	assert.Regexp(t, `\bNo traffic alarm so far\b`, s)
	assert.Regexp(t, `Traffic is normal`, s)
}
