package p4

import (
	"bytes"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

const (
	_group = "group"
	_user  = "user"
)

type Permission struct {
	Mode    string `json:"mode"`
	IsGroup bool   `json:"isGroup"`
	Name    string `json:"name"`
	Host    string `json:"host"`
	Path    string `json:"path"`
}

type ACL struct {
	List []*Permission
}

func (p *ACL) String() string {
	var (
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err := _protectionsTemplate.Parse(_protectionsTemplateTxt); err != nil {
		return ""
	}
	if err := _protectionsTemplate.Execute(contentBuf, p); err != nil {
		return ""
	}
	return contentBuf.String()
}

func newPermission(line string) *Permission {
	line = strings.TrimSpace(line)
	fields := strings.Split(line, " ")
	if len(fields) != 5 {
		return nil
	}
	isGroup := true
	if fields[1] == _user {
		isGroup = false
	}
	return &Permission{
		Mode:    fields[0],
		IsGroup: isGroup,
		Name:    fields[2],
		Host:    fields[3],
		Path:    fields[4],
	}
}

var (
	_protectionsTemplate = template.New("ACL config template")
)
var _protectionsTemplateTxt = `Protections:
{{- range .List }}
	{{.Mode}} {{if .IsGroup}} group {{else}} user {{end}} {{.Name}} {{.Host}} {{.Path}}
{{- end }}`

func (conn *Conn) WriteProtections(acl *ACL) (out []byte, err error) {
	if acl == nil {
		err = errors.New("Access control list is empty!")
		return
	}
	content := []byte(acl.String())
	return conn.Input([]string{"protect", "-i"}, content)
}

func (conn *Conn) Protections() (acl *ACL, err error) {
	var result []Result
	if result, err = conn.RunMarshaled("protect", []string{"-o"}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if acl, _ = result[0].(*ACL); acl == nil {
		err = errors.New("Type not match")
		return
	}
	return
}

func (conn *Conn) ProtectionsDump() (out []byte, err error) {
	return conn.Output([]string{"protect", "-o"})
}
