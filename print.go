package p4

import (
	"runtime"
)

func (conn *Conn) Print(path string) (content []byte, err error) {
	// -q : suppress header line.
	out, err := conn.Output([]string{"print", "-q", path})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (conn *Conn) Print2File(path string, outFile string) (err error) {
	if runtime.GOOS == "windows" {
		conn.env = append(conn.env, "P4CHARSET=cp936")
	}
	if _, err = conn.Output([]string{"print", "-q", "-o", outFile, path}); err != nil {
		return
	}
	return
}
