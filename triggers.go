package p4

import (
	"bytes"
	"text/template"
)

type Triggers struct {
	Lines []string
}

func (t *Triggers) String() string {
	var (
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err := _triggersTemplate.Parse(_triggerTemplateTxt); err != nil {
		return ""
	}
	if err := _triggersTemplate.Execute(contentBuf, t); err != nil {
		return ""
	}
	return contentBuf.String()
}

var (
	_triggersTemplate = template.New("Trigger config template")
)
var _triggerTemplateTxt = `Triggers:
{{- range .Lines }}
	{{.}}
{{- end }}`

func (conn *Conn) WriteTriggers(lines []string) (out []byte, err error) {
	t := &Triggers{Lines: lines}
	content := []byte(t.String())
	return conn.Input([]string{"triggers", "-i"}, content)
}

func (conn *Conn) Triggers() (lines []string, err error) {
	var (
		result   []Result
		triggers *Triggers
	)
	if result, err = conn.RunMarshaled("triggers", []string{"-o"}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if triggers, _ = result[0].(*Triggers); triggers == nil {
		return
	}
	lines = triggers.Lines
	return
}

func (conn *Conn) TriggersDump() (out []byte, err error) {
	return conn.Output([]string{"triggers", "-o"})
}
