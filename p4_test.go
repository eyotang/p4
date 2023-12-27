package p4

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func setup(t *testing.T) (*Conn, error) {
	address := os.Getenv("P4PORT") // ssl:localhost:1666
	user := os.Getenv("P4USER")
	password := os.Getenv("P4PASSWD")
	return NewConn(address, user, password)
}

func TestChange_OutputMaps(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test OutputMaps", t, func() {
		So(err, ShouldBeNil)

		Convey("Output Maps", func() {
			change, err := conn.OutputMaps("describe", "-S", "-s", "6534")
			So(change, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
