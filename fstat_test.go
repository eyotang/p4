package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFstat_FileExist(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test Fstat", t, func() {
		So(err, ShouldBeNil)

		Convey("File exist", func() {
			yes, err := conn.FileExist("//.swarm/triggers/create_swarm_review.py")
			So(yes, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("File not exist", func() {
			yes, err := conn.FileExist("//.swarm/triggers/create_swarm_review2.py")
			So(yes, ShouldBeFalse)
			So(err, ShouldBeNil)
		})

		Convey("File(dir) not exist", func() {
			yes, err := conn.FileExist("//.swarm/triggers/")
			So(yes, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
	})
}
