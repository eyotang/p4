package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Assumes sample depot is running on localhost:1666, and p4 binary is in path.
func TestDirs(t *testing.T) {
	var (
		conn *Conn
		err  error

		dirs  []*Dir
		files []*File
	)

	conn, err = setup(t)
	Convey("test Dirs functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List dirs", func() {
			dirs, err = conn.Dirs([]string{"//DM99.ZGame.Project/Main/ZGame_Mainline/*"})
			So(err, ShouldBeNil)
			So(len(dirs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("List files", func() {
			files, err = conn.Files([]string{"//DM99.ZGame.Project/Main/ZGame_Mainline/*"})
			So(err, ShouldBeNil)
			So(len(files), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
