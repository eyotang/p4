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

			conn.ChangeUser("guanxiao", "C72AFEA55FC4B855443D2429D2557BC8")
			message, err := conn.Reshelve(16309)
			fmt.Println(message)
			So(message, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestConn_DeleteShelve(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test changes", t, func() {
		So(err, ShouldBeNil)

		Convey("List shelved", func() {
			conn = conn.WithClient("root_arl_mainline")

			message, err := conn.DeleteShelve(17529)
			fmt.Println(message)
			So(message, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
