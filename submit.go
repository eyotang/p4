package p4

import (
	"strconv"
	"strings"
)

func (conn *Conn) SubmitShelve(change uint64, desc string) (message string, err error) {
	var (
		out []byte
	)
	args := []string{"submit", "-e", strconv.FormatUint(change, 10)}
	if len(desc) > 0 {
		args = append(args, []string{"-d", desc}...)
	}
	out, err = conn.Output(args)
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
