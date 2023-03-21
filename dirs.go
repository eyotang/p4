package p4

import "fmt"

type Dir struct {
	Dir string
}

func (f *Dir) String() string {
	return fmt.Sprintf("%s/", f.Dir)
}

func (conn *Conn) Dirs(paths []string) (dirs []*Dir, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("dirs", paths); err != nil {
		return
	}
	for idx := range results {
		if dir, ok := results[idx].(*Dir); !ok {
			continue
		} else {
			dirs = append(dirs, dir)
		}
	}
	return
}
