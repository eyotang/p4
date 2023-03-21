package p4

import (
	"os"
	"testing"
)

func setup(t *testing.T) (*Conn, error) {
	address := os.Getenv("P4PORT") // ssl:localhost:1666
	user := os.Getenv("P4USER")
	password := os.Getenv("P4PASSWD")
	return NewConn(address, user, password)
}
