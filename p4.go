// Copyright 2012 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p4

import (
	"bufio"
	"bytes"
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
)

// Conn is an interface to the Conn command line client.
type ConnOptions struct {
	address  string
	binary   string
	username string
	password string
}

type Conn struct {
	ConnOptions
	env []string
}

func NewConn(address, username, password string) (conn *Conn, err error) {
	conn = &Conn{
		ConnOptions: ConnOptions{
			binary:   "p4",
			address:  address,
			username: username,
			password: password,
		}}
	if err = conn.Login(); err != nil {
		return
	}
	return
}

type TagLine struct {
	Tag   string
	Value []byte
}

var tokenRegexp = regexp.MustCompile("([0-9A-Z]{32})")

func (p *Conn) Login() (err error) {
	env := []string{
		//"P4CLIENT=" + p.Client,
		"P4PORT=" + p.address,
		"P4USER=" + p.username,
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
		password bytes.Buffer
		token    bytes.Buffer
		errors   bytes.Buffer
	)

	password.Write([]byte(p.password))

	cmd := exec.Command(p.binary, "login", "-p")
	cmd.Env = env
	cmd.Stdin = &password
	cmd.Stdout = &token
	cmd.Stderr = &errors

	log.Println("running", cmd.Args)
	if err = cmd.Run(); err != nil {
		return P4Error{err, []string{"p4", "login"}, errors.Bytes()}
	}
	env = append(env, "P4PASSWD="+tokenRegexp.FindString(token.String()))
	p.env = env
	return
}

// Output runs p4 and captures stdout.
func (p *Conn) Output(args []string) ([]byte, error) {
	b := p.binary
	if !strings.Contains(b, "/") {
		b, _ = exec.LookPath(b)
	}
	cmd := exec.Cmd{
		Path: b,
		Args: []string{p.binary},
	}
	if p.env != nil {
		cmd.Env = p.env
	}
	if p.address != "" {
		cmd.Args = append(cmd.Args, "-p", p.address)
	}
	cmd.Args = append(cmd.Args, args...)

	//log.Println("running", cmd.Args)
	return cmd.Output()
}

// Runs p4 with -G and captures the result lines.
func (p *Conn) RunMarshaled(command string, args []string) (result []Result, err error) {
	var out []byte
	if out, err = p.Output(append([]string{"-G", command}, args...)); err != nil {
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
		cl, _ := strconv.ParseInt(r["change"], 10, 64)
		c.Change = int(cl)
		t, _ := strconv.ParseInt(r["time"], 10, 64)
		c.Time = int(t)
		return &c

	case "group":
		var users []string
		groupUserInfo := GroupUserInfo{
			Group: imap["Group"].(string),
		}
		for k, v := range imap {
			if strings.HasPrefix(k, "Users") {
				users = append(users, v.(string))
			}
		}
		groupUserInfo.Users = users
		return &groupUserInfo

	default:
		log.Panicf("unknown code %q", command)
	}
	return nil
}

func (p *Conn) Fstat(paths []string) (results []Result, err error) {
	r, err := p.RunMarshaled("fstat",
		append([]string{"-Of", "-Ol"}, paths...))
	return r, err
}

func (p *Conn) Files(paths []string) (results []Result, err error) {
	return p.RunMarshaled("files", paths)
}

func (p *Conn) Dirs(paths []string) ([]Result, error) {
	return p.RunMarshaled("dirs", paths)
}

func (p *Conn) Print(path string) (content []byte, err error) {
	// -q : suppress header line.
	out, err := p.Output([]string{"print", "-q", path})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (p *Conn) Changes(paths []string) ([]Result, error) {
	return p.RunMarshaled("changes", append([]string{"-l"}, paths...))
}

// P4Admin
func (p *Conn) Groups() (result []Result, err error) {
	var out []byte
	if out, err = p.Output([]string{"groups", "-i"}); err != nil {
		return
	}
	r := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		result = append(result, &GroupInfo{Group: scanner.Text()})
	}
	return
}

func (p *Conn) Members(group string) ([]Result, error) {
	if runtime.GOOS == "windows" {
		p.env = append(p.env, "P4CHARSET=cp936")
	}
	return p.RunMarshaled("group", []string{"-o", group})
}

////////////////
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

// Stat has the data for a single file revision.
type Stat struct {
	DepotFile   string
	HeadAction  string
	HeadType    string
	HeadTime    int64
	HeadRev     int64
	HeadChange  int64
	HeadModTime int64
	FileSize    int64
	Digest      string
}

func (f *Stat) String() string {
	return fmt.Sprintf("%s#%d - change %d (%s)",
		f.DepotFile, f.HeadRev, f.HeadChange, f.HeadType)
}

// File has the data for a single file.
type File struct {
	Code      string
	DepotFile string
	Revision  int64
	Action    string
	Type      string
	ModTime   int64
}

func (f *File) String() string {
	return f.DepotFile
}

type Dir struct {
	Dir string
}

func (f *Dir) String() string {
	return fmt.Sprintf("%s/", f.Dir)
}

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

type GroupInfo struct {
	Group string
}

func (g *GroupInfo) String() string {
	return fmt.Sprintf("group: %s", g.Group)
}

type GroupUserInfo struct {
	Group string
	Users []string
}

func (gu *GroupUserInfo) String() string {
	return fmt.Sprintf("group: %s, users: %v", gu.Group, gu.Users)
}
