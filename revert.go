package p4

import (
	"os/exec"
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

func (conn *Conn) ForceRevert(change uint64) (message string, err error) {
	var (
		out []byte
	)
	args := []string{"revert", "-f", "-c", strconv.FormatUint(change, 10), "//..."}
	out, err = exec.Command("p4", args...).CombinedOutput()
	message = strings.TrimSpace(string(out))
	return
}
