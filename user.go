package p4

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strings"
	"text/template"
)

type UserInfo struct {
	User       string `json:"user"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	AuthMethod string `json:"authMethod"`
}

func (u *UserInfo) String() string {
	buf, _ := json.Marshal(u)
	return string(buf)
}

func (conn *Conn) Users() (list []*UserInfo, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("users", []string{}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	for _, v := range result {
		if userInfo, ok := v.(*UserInfo); !ok {
			return
		} else {
			list = append(list, userInfo)
		}
	}
	return
}

func (conn *Conn) User(user string) (info *UserInfo, err error) {
	var (
		result []Result
	)
	if runtime.GOOS == "windows" {
		conn.env = append(conn.env, "P4CHARSET=cp936")
	}
	if result, err = conn.RunMarshaled("user", []string{"-o", user}); err != nil {
		return
	}
	if len(result) == 0 {
		return
	}
	if info, _ = result[0].(*UserInfo); info == nil {
		return
	}
	return
}

// DeleteUser
// https://www.perforce.com/manuals/cmdref/Content/CmdRef/p4_user.html#Examples
// ex: user=sammy
// 1. Delete sammy
// 2. Delete all of sammy's workspace clients, including those where a user other than sammy has files opened
// 3. Delete sammy from the protections table and groups
func (conn *Conn) DeleteUser(user string) (message string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"user", "-D", "-F", "-y", user}); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}

var _userTemplate = template.New("user template")
var _userTemplateTxt = `User:   {{ .User }}
Email:  {{ .Email }}
FullName:       {{ .FullName }}
AuthMethod:     ldap
`

// CreateUser 创建标准用户
func (conn *Conn) CreateUser(info *UserInfo) (message string, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
	)
	if info == nil {
		return
	}
	if _, err = _userTemplate.Parse(_userTemplateTxt); err != nil {
		return
	}
	if err = _userTemplate.Execute(contentBuf, info); err != nil {
		return
	}
	if out, err = conn.Input([]string{"user", "-i", "-f"}, contentBuf.Bytes()); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}
