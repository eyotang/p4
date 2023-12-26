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
	Change int
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
	Change      int
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
	_changeTemplate = template.New("ACL config template")
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
Files: {{-range .Files }}
        {{.}}
{{- end }}
`

func (conn *Conn) ChangeList(change int) (cl *ChangeList, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("change", append([]string{"-o", strconv.Itoa(change)})); err != nil {
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
