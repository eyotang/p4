package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func setup(t *testing.T) (*Conn, error) {
	//address := os.Getenv("P4PORT") // ssl:localhost:1666
	//user := os.Getenv("P4USER")
	//password := os.Getenv("P4PASSWD")

	address := "ssl:techcentertest.int.hypergryph.com:1666"
	user := "root"
	password := "p4super@2021"
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
			change, err := conn.OutputMaps("describe", "-S", "-s", "10067")
			So(change, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestConn_Trust(t *testing.T) {
	Convey("test Trust", t, func() {
		address := os.Getenv("P4PORT")
		msg, err := Trust(address)
		So(msg, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
