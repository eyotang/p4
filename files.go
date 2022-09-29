package p4

// File has the data for a single file.
type File struct {
	Code      string
	DepotFile string
	Revision  int64
	Action    string
	Type      string
	ModTime   int64
}

func (f *File) String() string {
	return f.DepotFile
}

func (conn *Conn) Files(paths []string) (results []Result, err error) {
	return conn.RunMarshaled("files", paths)
}
