package p4

import "log"

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

func (conn *Conn) Files(paths []string) (files []*File, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("files", paths); err != nil {
		return
	}
	for idx := range results {
		if file, ok := results[idx].(*File); !ok {
			log.Printf("type translate err: %s", results[idx])
			continue
		} else if file.Action == "delete" || file.Action == "move/delete" || file.Action == "purge" {
			continue
		} else {
			files = append(files, file)
		}
	}
	return
}
