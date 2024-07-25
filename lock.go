package p4

import (
	"strings"
)

func (conn *Conn) Unlock(file string) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"unlock", "-f", file})
	message = strings.TrimSpace(string(out))
	return
}
