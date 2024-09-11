package p4

import (
	"strconv"
	"strings"
)

func (conn *Conn) DeleteShelved(path string, change uint64) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"shelve", "-df", "-c", strconv.FormatUint(change, 10), path})
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) DeleteShelve(shelveCL uint64) (message string, err error) {
	var out []byte
	out, err = conn.Output([]string{"shelve", "-f", "-d", "-Af", "-c", strconv.FormatUint(shelveCL, 10)})
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) Reshelve(shelveCL uint64) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"reshelve", "-s", strconv.FormatUint(shelveCL, 10), "-f"})
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) Unshelve(shelveCL, CL uint64) (message string, err error) {
	var out []byte
	if CL != 0 {
		out, err = conn.Output([]string{"unshelve", "-s", strconv.FormatUint(shelveCL, 10), "-c", strconv.FormatUint(CL, 10), "-f"})
	} else {
		out, err = conn.Output([]string{"unshelve", "-s", strconv.FormatUint(shelveCL, 10), "-f"})
	}
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) UnshelveBypassExclusive(shelveCL, CL uint64) (message string, err error) {
	var out []byte
	if CL != 0 {
		out, err = conn.Output([]string{"unshelve", "--bypass-exclusive-lock", "-s", strconv.FormatUint(shelveCL, 10), "-c", strconv.FormatUint(CL, 10), "-f"})
	} else {
		out, err = conn.Output([]string{"unshelve", "--bypass-exclusive-lock", "-s", strconv.FormatUint(shelveCL, 10), "-c", strconv.FormatUint(CL, 10), "-f"})
	}
	message = strings.TrimSpace(string(out))
	return
}
