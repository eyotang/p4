package p4

import (
	"bytes"
	"fmt"
)

type P4Error struct {
	Status    error
	Arguments []string
	Output    []byte
}

func (err P4Error) Error() string {
	return fmt.Sprintf(
		"%s (%s): %s",
		err.Status,
		err.Arguments,
		string(bytes.TrimRight(err.Output, "\n")),
	)
}
