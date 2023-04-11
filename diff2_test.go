package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff2Change(t *testing.T) {
	var (
		conn *Conn
		err  error

		diffs []*Diff2
	)

	conn, err = setup(t)
	Convey("test Diff2 change functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Has update", func() {
			diffs, err = conn.Diff2Change("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 471, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 469)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("No change", func() {
			diffs, err = conn.Diff2Change("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 471, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 471)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("Not exist", func() {
			diffs, err = conn.Diff2Change("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 1, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 471)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDiff2Shelve(t *testing.T) {
	var (
		conn *Conn
		err  error

		diffs []*Diff2
	)

	conn, err = setup(t)
	Convey("test Diff2 shelve functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Has update", func() {
			diffs, err = conn.Diff2Shelve("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 472, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 473)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("No change", func() {
			diffs, err = conn.Diff2Shelve("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 473, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 473)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("Not exist", func() {
			diffs, err = conn.Diff2Shelve("//DM99.ZGame.Project/Main/ZGame_Mainline/...", 472, "//DM99.ZGame.Project/Main/ZGame_Mainline/...", 473)
			So(err, ShouldBeNil)
			So(len(diffs), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
