package p4

import (
	"strconv"
	"strings"
)

func (conn *Conn) DeleteShelved(path string, change int) (message string, err error) {
	var (
		out []byte
	)
	_, err = conn.Output([]string{"shelve", "-df", "-c", strconv.Itoa(change), path})
	message = strings.TrimSpace(string(out))
	return
}
