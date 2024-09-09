package p4

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConn_UserReshelve(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test changes", t, func() {
		So(err, ShouldBeNil)

		Convey("List shelved", func() {
			conn = conn.WithClient("root_arl_mainline")

			message, err := conn.UserReshelve(16184, "sunqi01")
			fmt.Println(message)
			So(message, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
