package p4

import (
	"fmt"
	"strings"
)

type Change struct {
	Desc   string
	User   string
	Status string
	Change int
	Time   int

	Path       string
	Code       string
	ChangeType string
	Client     string
}

func (c *Change) String() string {
	l := len(c.Desc)
	if l > 250 {
		l = 250
	}
	return fmt.Sprintf("change %d by %s - %s", c.Change, c.User, strings.Trim(c.Desc[:l], " "))
}

// Changes path格式: //Stream_Root/...
func (conn *Conn) Changes(paths []string) ([]Result, error) {
	return conn.RunMarshaled("changes", append([]string{"-l"}, paths...))
}

// Shelved path格式: //Stream_Root/...
func (conn *Conn) Shelved(path string) (shelved []*Change, err error) {
	var (
		result []Result
	)
	if result, err = conn.RunMarshaled("changes", append([]string{"-s", "shelved"}, path)); err != nil {
		return
	}
	for idx := range result {
		if r, ok := result[idx].(*Change); !ok {
			continue
		} else {
			shelved = append(shelved, r)
		}
	}
	return
}
