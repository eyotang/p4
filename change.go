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

func (conn *Conn) Changes(paths []string) ([]Result, error) {
	return conn.RunMarshaled("changes", append([]string{"-l"}, paths...))
}
