package ui

import (
	"bytes"
	"html/template"
	"path/filepath"
	"sync"
)

type Renderer interface {
	Render(*State) error
	Result() string
}

type renderer struct {
	sync.RWMutex
	template *template.Template
	result   string
}

func NewRenderer(baseDir string) (Renderer, error) {
	t, err := template.ParseGlob(filepath.Join(baseDir, "*.tmpl"))
	if err != nil {
		return nil, err
	}

	ret := &renderer{template: t}
	err = ret.Render(&State{})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *renderer) Render(s *State) error {
	const indexTemplateName = "index"

	buf := make([]byte, 0, 4096)
	wr := bytes.NewBuffer(buf)
	err := r.template.ExecuteTemplate(wr, indexTemplateName, s)
	if err != nil {
		return err
	}
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
