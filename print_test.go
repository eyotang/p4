package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConn_Print2File(t *testing.T) {
	conn, err := setup(t)
	Convey("test Print2File functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Print2File filename with chinese", func() {
			err = conn.Print2File("//DM99.ZGame.Project/Development/ZGame_ArtDev/Assets/Data/策划配表2.xlsx", "策划配表2.xlsx")
			So(err, ShouldBeNil)
		})
	})
}
