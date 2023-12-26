package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChange_Shelved(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test changes", t, func() {
		So(err, ShouldBeNil)

		Convey("List shelved", func() {
			shelved, err := conn.Shelved("//DM99.ZGame.Project/Development/test12/...")
			So(shelved, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("List empty stream shelved", func() {
			_, err := conn.Shelved("")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestChange_ChangeList(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test change", t, func() {
		So(err, ShouldBeNil)

		Convey("Display change", func() {
			change, err := conn.ChangeList(6534)
			So(change, ShouldNotBeNil)
			So(change.Type, ShouldEqual, "public")
			So(err, ShouldBeNil)
		})
	})
}
