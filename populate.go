package p4

import "strings"

func (conn *Conn) Populate(stream string) (message string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"populate", "-o", "-S", stream, "-r", "-d", "{1 more items}"}); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}
