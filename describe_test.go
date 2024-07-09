package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDescribeShelved(t *testing.T) {
	var (
		conn *Conn
		err  error

		description *Description
	)

	conn, err = setup(t)
	Convey("test DescribeShelved functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Describe Shelved", func() {
			description, err = conn.DescribeShelved(11941)
			So(err, ShouldBeNil)
			So(len(description.DepotFiles), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
