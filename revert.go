package p4

import (
	"strconv"
	"strings"
)

func (conn *Conn) Revert(change uint64) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"revert", "-c", strconv.FormatUint(change, 10), "//..."})
	message = strings.TrimSpace(string(out))
	return
}
