package p4

import "strings"

// Prune Removes unmodified files from a stream. 危险慎用
func (conn *Conn) Prune(stream string) (message string, err error) {
	var out []byte
	if out, err = conn.Output([]string{"prune", "-y", "-S", stream}); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}
