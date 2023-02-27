package p4

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

type StreamInfo struct {
	Stream  string
	Owner   string
	Name    string
	Parent  string // steam type为mainline时，parent必须为none，其余类型stream需要填写现有的stream（格式：//depotname/streamname）
	Type    string // mainline, development, release, virtual, task
	Options string // allsubmit unlocked notoparent nofromparent mergedown
}

func (s *StreamInfo) String() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (conn *Conn) Streams() (list []*StreamInfo, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("streams", []string{}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	for _, v := range result {
		if stream, ok := v.(*StreamInfo); !ok {
			return
		} else {
			list = append(list, stream)
		}
	}
	return
}

var _streamTemplate = template.New("stream template")
var _streamTemplateTxt = `Stream:  {{ .Stream }}
Owner:        {{ .Owner }}
Name:         {{ .Name }}
Parent:       {{ if eq .Type "mainline" }}none{{else}}{{ .Parent }}{{ end }}
Type:         {{ .Type }}
Description:
        Created by {{ .Owner }} automatically.
{{- if ne .Type "mainline" }}
Options:        {{ .Options }}
{{- end }}
Paths:
        share ...
`

func (conn *Conn) CreateStream(name, streamType, parent, location string, populate bool) (message string, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
		streamInfo = StreamInfo{
			Stream:  location,
			Owner:   conn.username,
			Name:    name,
			Parent:  parent,
			Type:    streamType,
			Options: "allsubmit unlocked toparent fromparent mergedown",
		}
	)
	if _, err = _streamTemplate.Parse(_streamTemplateTxt); err != nil {
		return
	}
	if err = _streamTemplate.Execute(contentBuf, streamInfo); err != nil {
		return
	}
	if out, err = conn.Input([]string{"stream", "-i"}, contentBuf.Bytes()); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	if populate {
		if _, err = conn.Populate(location); err != nil {
			return
		}
	}
	return
}

func (conn *Conn) DeleteStream(location string) (message string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"stream", "-d", location}); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}
