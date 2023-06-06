package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient_Clients(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test Clients", t, func() {
		So(err, ShouldBeNil)

		Convey("List clients", func() {
			clients, err := conn.Clients("//DM99.ZGame.Project/Main/ZGame_Mainline")
			So(clients, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
	})
}
