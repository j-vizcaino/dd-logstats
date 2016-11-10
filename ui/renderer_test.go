package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, s *State) string {
	r, err := NewRenderer("assets")
	assert.NoError(t, err)

	err = r.Render(s)
	assert.NoError(t, err)

	return r.Result()
}

func TestRendererAlarmOn(t *testing.T) {

	s := render(t, &State{
		AlarmIsActive: true,
		Alarms: []string{
			"alarm 1",
			"alarm 2",
		},
	})

	assert.Regexp(t, `\balarm 1\b`, s)
	assert.Regexp(t, `\balarm 2\b`, s)
	assert.Regexp(t, `High traffic`, s)
}

func TestRendererInitial(t *testing.T) {
	s := render(t, &State{})

	assert.Regexp(t, `\bNo traffic alarm so far\b`, s)
	assert.Regexp(t, `Traffic is normal`, s)
}
