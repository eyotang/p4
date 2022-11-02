package p4

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

type GroupUserInfo struct {
	Group     string
	Owners    []string
	SubGroups []string
	Users     []string
	Timestamp string
}

func (gu *GroupUserInfo) String() string {
	return fmt.Sprintf("group: %s, users: %v", gu.Group, gu.Users)
}

// P4Admin
// P4 用户组
func (conn *Conn) Groups() (result []string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"groups", "-i"}); err != nil {
		return
	}
	r := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return
}

// 用户属于的组，作为成员Member
func (conn *Conn) GroupsBelong(user string) (result []string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"groups", "-u", user}); err != nil {
		return
	}
	r := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return
}

// 用户拥有的组，作为拥有者Owner
func (conn *Conn) GroupsOwned(user string) (result []string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"groups", "-o", user}); err != nil {
		return
	}
	r := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return
}

func (conn *Conn) GroupUsers(group string) (members []string, err error) {
	var (
		result    []Result
		groupInfo *GroupUserInfo
	)
	if runtime.GOOS == "windows" {
		conn.env = append(conn.env, "P4CHARSET=cp936")
	}
	if result, err = conn.RunMarshaled("group", []string{"-o", group}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if groupInfo, _ = result[0].(*GroupUserInfo); groupInfo == nil {
		return
	}
	members = groupInfo.Users
	return
}
func (conn *Conn) GroupSubGroups(group string) (subGroups []string, err error) {
	var (
		result    []Result
		groupInfo *GroupUserInfo
	)
	if runtime.GOOS == "windows" {
		conn.env = append(conn.env, "P4CHARSET=cp936")
	}
	if result, err = conn.RunMarshaled("group", []string{"-o", group}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if groupInfo, _ = result[0].(*GroupUserInfo); groupInfo == nil {
		return
	}
	subGroups = groupInfo.SubGroups
	return
}

func (conn *Conn) ExistGroup(group string) (yes bool, err error) {
	var (
		groups []string
	)
	if groups, err = conn.Groups(); err != nil {
		return
	}
	for _, v := range groups {
		if v == group {
			yes = true
			return
		}
	}
	return
}

var _groupTemplate = template.New("group template")
var _groupTemplateTxt = `Group:  {{ .Group }}
Description:
	Auto generated at {{ .Timestamp }}.
MaxResults:     unset
MaxScanRows:    unset
MaxLockTime:    unset
MaxOpenFiles:   unset
Timeout:        unlimited
PasswordTimeout:        unlimited
Subgroups:
{{- range .SubGroups }}
	{{.}}
{{- end }}
Owners:
{{- range .Owners }}
	{{.}}
{{- end }}
Users:
{{- range .Users }}
	{{.}}
{{- end }}
`

// 需要较高权限
func (conn *Conn) CreateGroup(group string, owners, subGroups, members []string) (message string, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
		groupInfo  = GroupUserInfo{
			Group:     group,
			Owners:    owners,
			SubGroups: subGroups,
			Users:     members,
			Timestamp: time.Now().Format("2006-01-02_15-04-05"),
		}
	)
	if _, err = _groupTemplate.Parse(_groupTemplateTxt); err != nil {
		return
	}
	if err = _groupTemplate.Execute(contentBuf, groupInfo); err != nil {
		return
	}
	if out, err = conn.Input([]string{"group", "-i"}, contentBuf.Bytes()); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) DeleteGroup(group string) (message string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"group", "-d", group}); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) AddGroupUsers(group string, owners, addMembers []string) (message string, err error) {
	var (
		yes       bool
		members   []string
		subGroups []string
	)
	if yes, err = conn.ExistGroup(group); err != nil {
		return
	} else if !yes {
		err = errors.Errorf("Group '%s' isn't exist!", group)
		return
	} else {
		if members, err = conn.GroupUsers(group); err != nil {
			return
		}
		if subGroups, err = conn.GroupSubGroups(group); err != nil {
			return
		}
	}

	members = append(members, addMembers...)
	return conn.CreateGroup(group, owners, subGroups, members)
}

func (conn *Conn) RemoveGroupUsers(group string, owners, removeMembers []string) (message string, err error) {
	var (
		yes                 bool
		members, newMembers []string
		subGroups           []string
	)
	if yes, err = conn.ExistGroup(group); err != nil {
		return
	} else if !yes {
		message = fmt.Sprintf("Group '%s' is not exist!", group)
		return
	}
	if members, err = conn.GroupUsers(group); err != nil {
		return
	}
	if subGroups, err = conn.GroupSubGroups(group); err != nil {
		return
	}
	for _, v := range members {
		shouldRemove := false
		for _, remove := range removeMembers {
			if v == remove {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newMembers = append(newMembers, v)
		}
	}
	return conn.CreateGroup(group, owners, subGroups, newMembers)
}
