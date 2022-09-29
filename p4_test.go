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

// Assumes sample depot is running on localhost:1666, and p4 binary is in path.
func TestDirs(t *testing.T) {
	var (
		address  = "localhost:1666"
		user     = "tangyongqiang"
		password = "123456"
	)

	c, err := NewConn(address, user, password)
	if err != nil {
		t.Fatalf("NewConn failed! err: %+v", err)
		return
	}
	rs, err := c.Dirs([]string{"//depot/*@700"})
	if err != nil {
		t.Fatalf("p4.Dirs: %v", err)
	}

	if len(rs) != 1 {
		t.Fatalf("p4.Dirs got: %v, want 1 result", rs)
	}

	d := rs[0].(*Dir)
	if d.Dir != "//depot/Jam" {
		t.Fatalf("p4.Dirs got dir %q", d.Dir)
	}
}
