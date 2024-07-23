package p4

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"text/template"
)

type Change struct {
	Desc   string
	User   string
	Status string
	Change uint64
	Time   int

	Path       string
	Code       string
	ChangeType string
	Client     string
}

func (c *Change) String() string {
	l := len(c.Desc)
	if l > 250 {
		l = 250
	}
	return fmt.Sprintf("change %d by %s - %s", c.Change, c.User, strings.Trim(c.Desc[:l], " "))
}

// Changes path格式: //Stream_Root/...
func (conn *Conn) Changes(paths []string) ([]Result, error) {
	return conn.RunMarshaled("changes", append([]string{"-l"}, paths...))
}

// Shelved path格式: //Stream_Root/...
func (conn *Conn) Shelved(path string) (shelved []*Change, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("changes", append([]string{"-s", "shelved"}, path)); err != nil {
		return
	}
	for idx := range result {
		if r, ok := result[idx].(*Change); !ok {
			continue
		} else {
			shelved = append(shelved, r)
		}
	}
	return
}

type ChangeList struct {
	Change      uint64
	Date        string
	Client      string
	User        string
	Status      string
	Type        string
	Description string
	ImportedBy  string
	Identity    string
	Jobs        []string
	Stream      string
	Files       []string
}

type NewChangeList struct {
	Change      string // 固定传"new"
	Date        string
	Client      string
	User        string
	Status      string
	Type        string
	Description string
	ImportedBy  string
	Identity    string
	Jobs        []string
	Stream      string
	Files       []string
}

func (cl *ChangeList) String() string {
	var (
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err := _changeTemplate.Parse(_changeTemplateTxt); err != nil {
		return ""
	}
	if err := _changeTemplate.Execute(contentBuf, cl); err != nil {
		return ""
	}
	return contentBuf.String()
}

var (
	_changeTemplate = template.New("change template")
)
var _changeTemplateTxt = `Change: {{.Change}}
Date: {{.Date}}
Client: {{.Client}}
User: {{.User}}
Status: {{.Status}}
Type: {{.Type}}
Description: {{.Description}}
ImportedBy: {{.ImportedBy}}
Identity: {{.Identity}}
Jobs: {{- range .Jobs }}
        {{.}}
{{- end }}
Stream: {{.Stream}}
Files: {{- range .Files }}
        {{.}}
{{- end }}
`

func (conn *Conn) ChangeList(change uint64) (cl *ChangeList, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("change", append([]string{"-o", strconv.FormatUint(change, 10)})); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if cl, _ = result[0].(*ChangeList); cl == nil {
		err = errors.New("Type not match")
		return
	}
	return
}

func (conn *Conn) UpdateChangeList(cl ChangeList) (message string, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err = _changeTemplate.Parse(_changeTemplateTxt); err != nil {
		return
	}
	if err = _changeTemplate.Execute(contentBuf, cl); err != nil {
		return
	}
	if out, err = conn.Input([]string{"change", "-f", "-i"}, contentBuf.Bytes()); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) ChangeListStream(change uint64) (stream string, err error) {
	var (
		changeList *ChangeList
		client     *Client
	)
	// 通过CL号，拿到Client（workspace）
	if changeList, err = conn.ChangeList(change); err != nil {
		return
	}
	if changeList == nil {
		err = errors.Errorf("change '%d' can NOT found", change)
		return
	}

	// 查看Client（workspace）里面配置的stream
	if client, err = conn.Client(changeList.Client); err != nil {
		return
	}
	if client == nil {
		err = errors.Errorf("client '%s' for change '%d' can NOT found", changeList.Client, change)
		return
	}

	stream = client.Stream

	return
}

func (conn *Conn) NewChangeList(cl NewChangeList) (change uint64, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err = _changeTemplate.Parse(_changeTemplateTxt); err != nil {
		return
	}
	if err = _changeTemplate.Execute(contentBuf, cl); err != nil {
		return
	}
	if out, err = conn.Input([]string{"change", "-i"}, contentBuf.Bytes()); err != nil {
		return
	}
	message := strings.TrimSpace(string(out))
	if strings.HasPrefix(message, "Change") {
		_, err = fmt.Sscanf(message, "Change %d created.", &change)
		if err != nil {
			return 0, fmt.Errorf("failed to parse change list number: %v", err)
		}
	}
	return
}
