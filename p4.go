// Copyright 2012 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p4

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ConnOptions Conn is an interface to the Conn command line client.
type ConnOptions struct {
	address  string
	binary   string
	username string
	password string
	client   string
}

type Conn struct {
	ConnOptions
	env []string
}

func NewConn(address, username, password string) (conn *Conn, err error) {
	return NewClientConn(address, username, password, "")
}

func NewClientConn(address, username, password, client string) (conn *Conn, err error) {
	conn = &Conn{
		ConnOptions: ConnOptions{
			binary:   "p4",
			address:  address,
			username: username,
			password: password,
		}}
	if client != "" {
		conn.client = client
	}
	if err = conn.Login(); err != nil {
		return
	}
	return
}

var tokenRegexp = regexp.MustCompile("([0-9A-Z]{32})")

func (conn *Conn) Login() (err error) {
	env := []string{
		"P4PORT=" + conn.address,
		"P4USER=" + conn.username,
	}
	if conn.client != "" {
		env = append(env, "P4CLIENT="+conn.client)
	}
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		env = append(env, "P4TRUST="+path.Join(home, "p4trust.txt"))
	} else {
		home := os.Getenv("HOME")
		env = append(env, "P4TRUST="+path.Join(home, ".p4trust"))
	}
	//fmt.Println(env)

	var (
		password = bytes.NewBufferString(conn.password)
		token    bytes.Buffer
		stderr   bytes.Buffer
	)

	cmd := exec.Command(conn.binary, "login", "-p")
	cmd.Env = env
	cmd.Stdin = password
	cmd.Stdout = &token
	cmd.Stderr = &stderr

	log.Println("running", cmd.Args)
	if err = cmd.Run(); err != nil {
		return P4Error{err, []string{"p4", "login"}, stderr.Bytes()}
	}
	env = append(env, "P4PASSWD="+tokenRegexp.FindString(token.String()))
	conn.env = env
	return
}

func (conn *Conn) WithClient(client string) *Conn {
	if client != "" {
		conn.env = append(conn.env, "P4CLIENT="+client)
	}
	return conn
}

// Output runs p4 and captures stdout.
func (conn *Conn) Output(args []string) (out []byte, err error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	b := conn.binary
	if !strings.Contains(b, "/") {
		b, _ = exec.LookPath(b)
	}
	cmd := exec.Cmd{
		Path:   b,
		Args:   []string{conn.binary},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if conn.env != nil {
		cmd.Env = conn.env
	}
	if conn.address != "" {
		cmd.Args = append(cmd.Args, "-p", conn.address)
	}
	cmd.Args = append(cmd.Args, args...)

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err, stderr.String())
	}
	out = stdout.Bytes()
	return
}

func (conn *Conn) Input(args []string, input []byte) (out []byte, err error) {
	var (
		content = bytes.NewBuffer(input)
		stdout  bytes.Buffer
		stderr  bytes.Buffer
	)
	b := conn.binary
	if !strings.Contains(b, "/") {
		b, _ = exec.LookPath(b)
	}
	cmd := exec.Cmd{
		Path:   b,
		Args:   []string{conn.binary},
		Stdin:  content,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if conn.env != nil {
		cmd.Env = conn.env
	}
	if conn.address != "" {
		cmd.Args = append(cmd.Args, "-p", conn.address)
	}
	cmd.Args = append(cmd.Args, args...)

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err, stderr.String())
	}
	out = stdout.Bytes()
	return
}

var (
	JSONArgs = []string{"-Mj", "-ztag"}
)

func (conn *Conn) OutputMaps(args ...string) (result []map[string]string, err error) {
	var (
		line   []byte
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	b := conn.binary
	if !strings.Contains(b, "/") {
		b, _ = exec.LookPath(b)
	}
	cmd := exec.Cmd{
		Path:   b,
		Args:   []string{conn.binary},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if conn.env != nil {
		cmd.Env = conn.env
	}
	if conn.address != "" {
		cmd.Args = append(cmd.Args, "-p", conn.address)
	}
	cmd.Args = append(cmd.Args, JSONArgs...)
	cmd.Args = append(cmd.Args, args...)

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err, stderr.String())
	}

	result = make([]map[string]string, 0)
	reader := bufio.NewReader(&stdout)
	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(line) <= 0 {
			continue
		}
		r := make(map[string]string)
		if err = json.Unmarshal(line, &r); err != nil {
			return
		}
		result = append(result, r)
	}
	return
}

// RunMarshaled p4 with -G and captures the result lines.
func (conn *Conn) RunMarshaled(command string, args []string) (result []Result, err error) {
	var (
		out []byte
	)
	if out, err = conn.Output(append([]string{"-G", command}, args...)); err != nil {
		return
	}
	r := bytes.NewBuffer(out)
	for {
		v, err := Decode(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		asMap, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("format err: p4 marshaled %v", v)
		}
		result = append(result, interpretResult(asMap, command))
	}

	if len(result) > 0 {
		err = nil
	}

	return result, err
}

func interpretResult(in map[interface{}]interface{}, command string) Result {
	imap := map[string]interface{}{}
	for k, v := range in {
		imap[k.(string)] = v
	}
	code := imap["code"].(string)
	if code == "error" {
		e := Error{}
		e.Severity = int(imap["severity"].(int32))
		e.Generic = int(imap["generic"].(int32))
		e.Data = imap["data"].(string)
		return &e
	}

	switch command {
	case "dirs":
		return &Dir{Dir: imap["dir"].(string)}
	case "files":
		r := map[string]string{}
		for k, v := range imap {
			r[k] = v.(string)
		}
		f := File{
			Code:      r["code"],
			DepotFile: r["depotFile"],
			Action:    r["action"],
			Type:      r["type"],
		}
		f.Revision, _ = strconv.ParseInt(r["rev"], 10, 64)
		f.ModTime, _ = strconv.ParseInt(r["time"], 10, 64)
		return &f
	case "fstat":
		r := map[string]string{}
		for k, v := range imap {
			r[k] = v.(string)
		}

		st := Stat{
			DepotFile:  r["depotFile"],
			HeadAction: r["headAction"],
			Digest:     r["digest"],
			HeadType:   r["headType"],
		}

		// Brilliant. We get the integers as decimal strings. Sigh.
		st.HeadTime, _ = strconv.ParseInt(r["headTime"], 10, 64)
		st.HeadRev, _ = strconv.ParseInt(r["headRev"], 10, 64)
		st.HeadChange, _ = strconv.ParseInt(r["headChange"], 10, 64)
		st.HeadModTime, _ = strconv.ParseInt(r["headModTime"], 10, 64)
		st.FileSize, _ = strconv.ParseInt(r["fileSize"], 10, 64)
		return &st

	case "changes":
		r := map[string]string{}
		for k, v := range imap {
			r[k] = v.(string)
		}
		c := Change{
			Desc:       r["desc"],
			User:       r["user"],
			Status:     r["status"],
			Path:       r["path"],
			Code:       r["code"],
			ChangeType: r["changeType"],
			Client:     r["client"],
		}
		cl, _ := strconv.ParseUint(r["change"], 10, 64)
		c.Change = cl
		t, _ := strconv.ParseInt(r["time"], 10, 64)
		c.Time = int(t)
		return &c

	case "change":
		var (
			stream string
		)
		r := map[string]string{}
		for k, v := range imap {
			r[k] = v.(string)
		}
		if v, exist := imap["Stream"]; exist {
			stream = v.(string)
		}
		cl := ChangeList{
			Date:        r["Date"],
			Client:      r["Client"],
			User:        r["User"],
			Status:      r["Status"],
			Type:        r["Type"],
			Description: r["Description"],
			ImportedBy:  r["ImportedBy"],
			Identity:    r["Identity"],
			Stream:      stream,
		}
		cl.Change, _ = strconv.ParseUint(r["Change"], 10, 64)
		return &cl

	case "group":
		var (
			owners    []string
			users     []string
			subGroups []string
		)
		groupUserInfo := GroupInfo{
			Group: imap["Group"].(string),
		}
		for k, v := range imap {
			if strings.HasPrefix(k, "Users") {
				users = append(users, v.(string))
			} else if strings.HasPrefix(k, "Owners") {
				owners = append(owners, v.(string))
			} else if strings.HasPrefix(k, "Subgroups") {
				subGroups = append(subGroups, v.(string))
			}
		}
		groupUserInfo.Owners = owners
		groupUserInfo.Users = users
		groupUserInfo.SubGroups = subGroups
		return &groupUserInfo

	case "triggers":
		var triggers Triggers
		for k, v := range imap {
			if strings.HasPrefix(k, "Triggers") {
				triggers.Lines = append(triggers.Lines, v.(string))
			}
		}
		return &triggers

	case "protect":
		var (
			ok         bool
			err        error
			idx        int
			permission *Permission
			acl        = ACL{
				store: make(map[int]*Permission),
			}
		)
		for k, v := range imap {
			if strings.HasPrefix(k, "ProtectionsComment") && len(v.(string)) > 0 {
				suffix := strings.TrimPrefix(k, "ProtectionsComment")
				if idx, err = strconv.Atoi(suffix); err != nil {
					continue
				}
				if permission, ok = acl.store[idx]; !ok {
					if err, permission = newComment(v.(string)); err != nil {
						continue
					}
					acl.store[idx] = permission
				} else {
					permission.updateComment(v.(string))
				}
			} else if strings.HasPrefix(k, "Protections") && len(v.(string)) > 0 {
				suffix := strings.TrimPrefix(k, "Protections")
				if idx, err = strconv.Atoi(suffix); err != nil {
					continue
				}
				if permission, ok = acl.store[idx]; !ok {
					if err, permission = newPermission(v.(string)); err != nil {
						continue
					}
					acl.store[idx] = permission
				} else {
					permission.updatePermit(v.(string))
				}
			}
		}
		return &acl

	case "stream", "streams":
		stream := StreamInfo{
			Stream:  imap["Stream"].(string),
			Owner:   imap["Owner"].(string),
			Name:    imap["Name"].(string),
			Parent:  imap["Parent"].(string),
			Type:    imap["Type"].(string),
			Options: imap["Options"].(string),
		}
		return &stream

	case "diff2":
		r := map[string]string{}
		for k, v := range imap {
			r[k] = v.(string)
		}
		v1, _ := strconv.ParseUint(r["rev"], 10, 64)
		v2, _ := strconv.ParseUint(r["rev2"], 10, 64)
		diff := Diff2{
			Code:   r["code"],
			Status: r["status"],
			DiffFile1: &DiffFile{
				DepotFile: r["depotFile"],
				Revision:  v1,
				Type:      r["type"],
			},
			DiffFile2: &DiffFile{
				DepotFile: r["depotFile2"],
				Revision:  v2,
				Type:      r["type2"],
			},
		}
		return &diff

	case "clients":
		var host string
		if v, exist := imap["Host"]; exist {
			host = v.(string)
		}
		client := Client{
			Owner:       imap["Owner"].(string),
			Client:      imap["client"].(string),
			Root:        imap["Root"].(string),
			Host:        host,
			Stream:      imap["Stream"].(string),
			Description: imap["Description"].(string),
		}
		return &client

	case "client":
		var (
			views  []string
			stream string
			host   string
		)
		for k, v := range imap {
			if strings.HasPrefix(k, "View") {
				views = append(views, v.(string))
			}
		}
		if v, exist := imap["Stream"]; exist {
			stream = v.(string)
		}
		if v, exist := imap["Host"]; exist {
			host = v.(string)
		}
		client := Client{
			Client:        imap["Client"].(string),
			Owner:         imap["Owner"].(string),
			Host:          host,
			Description:   strings.TrimSpace(imap["Description"].(string)),
			Root:          imap["Root"].(string),
			Options:       imap["Options"].(string),
			SubmitOptions: imap["SubmitOptions"].(string),
			Stream:        stream,
			Type:          imap["Type"].(string),
			View:          views,
		}
		return &client

	case "user", "users":
		var (
			authMethod string
		)
		if v, exist := imap["AuthMethod"]; exist {
			authMethod = v.(string)
		}
		user := UserInfo{
			User:       imap["User"].(string),
			Email:      imap["Email"].(string),
			FullName:   imap["FullName"].(string),
			AuthMethod: authMethod,
		}
		return &user

	case "describe":
		var path string
		if v, ok := imap["path"]; ok {
			path = v.(string)
		}
		d := Description{
			Change:     imap["change"].(string),
			User:       imap["user"].(string),
			Describe:   imap["desc"].(string),
			ChangeType: imap["changeType"].(string),
			Path:       path,
			Time:       imap["time"].(string),
			Client:     imap["client"].(string),
			Status:     imap["status"].(string),
		}
		return &d

	default:
		log.Panicf("unknown command %q", command)
	}
	return nil
}

// Result //////////////
type Result interface {
	String() string
}

type Error struct {
	Generic  int
	Severity int
	Data     string
}

func (e *Error) String() string {
	return fmt.Sprintf("error %d(%d): %s", e.Generic, e.Severity, e.Data)
}

func (e *Error) Error() string {
	return e.String()
}
