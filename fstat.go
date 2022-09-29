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
}

func (f *Stat) String() string {
	return fmt.Sprintf("%s#%d - change %d (%s)",
		f.DepotFile, f.HeadRev, f.HeadChange, f.HeadType)
}

func (conn *Conn) Fstat(paths []string) (results []Result, err error) {
	r, err := conn.RunMarshaled("fstat",
		append([]string{"-Of", "-Ol"}, paths...))
	return r, err
}
