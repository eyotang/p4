package p4

import "fmt"

// Stat has the data for a single file revision.
type Stat struct {
	DepotFile   string
	HeadAction  string
	HeadType    string
	HeadTime    int64
	HeadRev     int64
	HeadChange  int64
	HeadModTime int64
	FileSize    int64
	Digest      string
	OtherLock   string
}

func (f *Stat) String() string {
	return fmt.Sprintf("%s#%d - change %d (%s)",
		f.DepotFile, f.HeadRev, f.HeadChange, f.HeadType)
}

func (conn *Conn) Fstat(paths []string) (results []Result, err error) {
	r, err := conn.RunMarshaled("fstat",
		append([]string{"-Of", "-Olh"}, paths...))
	return r, err
}

func (conn *Conn) Fstats(paths []string) (stats []*Stat, err error) {
	results, err := conn.RunMarshaled("fstat", paths)
	for _, result := range results {
		if stat, ok := result.(*Stat); !ok {
			continue
		} else {
			stats = append(stats, stat)
		}
	}
	return
}

func (conn *Conn) FileExist(path string) (yes bool, err error) {
	var result []Result
	if result, err = conn.RunMarshaled("fstat", []string{path}); err != nil {
		return
	}

	if len(result) > 0 {
		_, isError := result[0].(*Error)
		yes = !isError
		return
	}

	return
}
