package p4

import "fmt"

type Dir struct {
	Dir string
}

func (f *Dir) String() string {
	return fmt.Sprintf("%s/", f.Dir)
}

func (conn *Conn) Dirs(paths []string) ([]Result, error) {
	return conn.RunMarshaled("dirs", paths)
}
