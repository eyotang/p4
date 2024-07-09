package p4

import (
	"log"
	"strconv"
)

// Description has the describe result of a single changelist.
type Description struct {
	Change     string
	User       string
	Describe   string
	ChangeType string
	Path       string
	Client     string
	Time       string
	Status     string
	DepotFiles []*File
}

func (d *Description) String() string {
	return d.Describe
}

func (conn *Conn) Describe(number uint64) (desc *Description, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("describe", []string{strconv.FormatUint(number, 10)}); err != nil {
		return
	}
	for idx := range results {
		if d, ok := results[idx].(*Description); !ok {
			log.Printf("type translate err: %s", results[idx])
			continue
		} else {
			// Get the first valid describe result and return it.
			desc = d
			return
		}
	}
	return
}

func (conn *Conn) DescribeShelved(number uint64) (desc *Description, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("describe", []string{"-S", "-s", strconv.FormatUint(number, 10)}); err != nil {
		return
	}
	for idx := range results {
		if d, ok := results[idx].(*Description); !ok {
			log.Printf("type translate err: %s", results[idx])
			continue
		} else {
			// Get the first valid describe result and return it.
			desc = d
			return
		}
	}
	return
}
