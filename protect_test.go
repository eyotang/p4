package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProtections_Protections(t *testing.T) {
	var (
		currentACL *ACL
		newACL     = &ACL{List: []*Permission{
			{"super", false, "root", "*", "//...", "更新于: 2022-10-09 11:17, 更新人: tangyongqiang, 描述: 评审控制"},
			{"write", true, "eyotang", "*", "//main/...", ""},
		}}
	)
	conn, err := setup(t)
	Convey("test Protections functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List triggers", func() {
			currentACL, err = conn.Protections()
			So(err, ShouldBeNil)
			So(len(currentACL.List), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("Preview permissions", func() {
			content := newACL.String()
			So(content, ShouldNotBeEmpty)
		})

		Convey("Write new permissions", func() {
			_, err = conn.WriteProtections(newACL)
			So(err, ShouldBeNil)
		})

		Convey("Restore old permissions", func() {
			_, err = conn.WriteProtections(currentACL)
			So(err, ShouldBeNil)
		})
	})
}
