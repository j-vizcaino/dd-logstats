package ui

import (
	"bytes"
	"html/template"
	"path/filepath"
	"sync"
	"time"
)

type Renderer interface {
	Render(*State) error
	Result() string
}

type renderer struct {
	sync.RWMutex
	template *template.Template
	result   string
	header   string
	footer   string
}

func renderString(t *template.Template, name string, data interface{}) (string, error) {
	buf := make([]byte, 0, 1024)
	wr := bytes.NewBuffer(buf)

	err := t.ExecuteTemplate(wr, name, data)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

func NewRenderer(
	baseDir string,
	statsPeriod time.Duration,
	alarmThreshold uint64,
	alarmPeriod time.Duration,
) (Renderer, error) {

	t, err := template.ParseGlob(filepath.Join(baseDir, "*.gohtml"))
	if err != nil {
		return nil, err
	}

	var header, footer string
	cfg := struct {
		RefreshPeriod  uint64
		AlarmPeriod    time.Duration
		AlarmThreshold uint64
		StatsPeriod    time.Duration
	}{
		// Shannon: sample twice as fast
		RefreshPeriod:  uint64(statsPeriod.Seconds() / 2),
		AlarmPeriod:    alarmPeriod,
		AlarmThreshold: alarmThreshold,
		StatsPeriod:    statsPeriod,
	}

	header, err = renderString(t, "header", cfg)
	if err != nil {
		return nil, err
	}
	footer, err = renderString(t, "footer", cfg)
	if err != nil {
		return nil, err
	}
	ret := &renderer{
		template: t,
		header:   header,
		footer:   footer,
	}
	err = ret.Render(&State{})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *renderer) Render(s *State) error {
	// No need to lock when using header and footer because
	// they are const and never touched after initialization
	wr := bytes.NewBufferString(r.header)
	err := r.template.ExecuteTemplate(wr, "body", s)
	if err != nil {
		return err
	}
	wr.WriteString(r.footer)
	r.Lock()
	defer r.Unlock()
	r.result = wr.String()
	return nil
}

func (r *renderer) Result() string {
	r.RLock()
	defer r.RUnlock()
	return r.result
}
