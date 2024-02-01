package p4

import (
	"strconv"
	"strings"
)

func (conn *Conn) SubmitShelve(change uint64) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"submit", "-e", strconv.FormatUint(change, 10)})
	message = strings.TrimSpace(string(out))
	return
}

func (conn *Conn) SubmitChange(change uint64) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"submit", "-c", strconv.FormatUint(change, 10)})
	message = strings.TrimSpace(string(out))
	return
}
