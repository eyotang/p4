package p4

import (
	"bytes"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
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
	Comment string `json:"comment"`
}

type ACL struct {
	store map[int]*Permission
	List  []*Permission
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

func newPermission(line string) (err error, p *Permission) {
	p = new(Permission)
	if err = p.updatePermit(line); err != nil {
		return
	}
	return
}

func newComment(line string) (err error, p *Permission) {
	p = new(Permission)
	if err = p.updateComment(line); err != nil {
		return
	}
	return
}

func (permission *Permission) updatePermit(line string) (err error) {
	line = strings.TrimSpace(line)
	fields := strings.Split(line, " ")
	if len(fields) != 5 {
		err = errors.New("Invalid format")
		return
	}
	isGroup := true
	if fields[1] == _user {
		isGroup = false
	}
	permission.Mode = fields[0]
	permission.IsGroup = isGroup
	permission.Name = fields[2]
	permission.Host = fields[3]
	permission.Path = fields[4]
	return
}

func (permission *Permission) updateComment(line string) (err error) {
	if len(line) <= 0 {
		err = errors.New("Comment is empty")
		return
	}
	comment := strings.TrimSpace(line)
	if !strings.HasPrefix(comment, "##") {
		err = errors.New("Comment format invalid")
		return
	}
	comment = strings.TrimPrefix(comment, "##")
	comment = strings.TrimSpace(comment)
	permission.Comment = comment
	return
}

var (
	_protectionsTemplate = template.New("ACL config template")
)
var _protectionsTemplateTxt = `Protections:
{{- range .List }}
	{{.Mode}} {{if .IsGroup}} group {{else}} user {{end}} {{.Name}} {{.Host}} {{.Path}}{{if .Comment}} ## {{.Comment}} {{end}}
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
	var (
		keys   []int
		result []Result
	)
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

	for k := range acl.store {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		acl.List = append(acl.List, acl.store[k])
	}
	return
}

func (conn *Conn) ProtectionsDump() (out []byte, err error) {
	return conn.Output([]string{"protect", "-o"})
}
