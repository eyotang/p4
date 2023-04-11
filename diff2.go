package p4

import (
	"encoding/json"
	"strconv"
)

type DiffFile struct {
	DepotFile string `json:"depotFile"`
	Revision  uint64 `json:"revision"`
	Type      string `json:"type"` // binary+l
}

type Diff2 struct {
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	DiffFile1 *DiffFile `json:"diffFile1"`
	DiffFile2 *DiffFile `json:"diffFile2"`
}

func (diff *Diff2) String() string {
	buf, _ := json.Marshal(diff)
	return string(buf)
}

// Diff2 参考手册：
// https://www.perforce.com/manuals/cmdref/Content/CmdRef/p4_diff2.html#p4_diff2
func (conn *Conn) Diff2(myStreamSpec, yourStreamSpec string) (diffs []*Diff2, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("diff2", []string{myStreamSpec, yourStreamSpec}); err != nil {
		return
	}
	for idx := range results {
		if diff, ok := results[idx].(*Diff2); !ok {
			continue
		} else {
			diffs = append(diffs, diff)
		}
	}
	return
}

func (conn *Conn) Diff2Change(myStream string, myChange uint, yourStream string, yourChange uint) ([]*Diff2, error) {
	return conn.Diff2(myStream+"@"+strconv.FormatUint(uint64(myChange), 10), yourStream+"@"+strconv.FormatUint(uint64(yourChange), 10))
}

func (conn *Conn) Diff2Shelve(myStream string, myShelve uint, yourStream string, yourShelve uint) ([]*Diff2, error) {
	return conn.Diff2(myStream+"@="+strconv.FormatUint(uint64(myShelve), 10), yourStream+"@="+strconv.FormatUint(uint64(yourShelve), 10))
}
