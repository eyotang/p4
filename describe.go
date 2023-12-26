package p4

import "log"

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
}

func (d *Description) String() string {
	return d.Describe
}

func (conn *Conn) Describe(number string) (desc *Description, err error) {
	var (
		results []Result
	)
	if results, err = conn.RunMarshaled("describe", []string{number}); err != nil {
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
