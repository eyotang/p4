package p4

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroup_Groups(t *testing.T) {
	var (
		message string
		group   = "group-xxx"
		users   = []string{"eyotang", "tangyongqiang"}
		other   = []string{"abc"}
		owners  = []string{"owner"}
	)
	conn, err := setup(t)
	Convey("test Group functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List groups", func() {
			var (
				groups []string
			)

			groups, err = conn.Groups()
			So(err, ShouldBeNil)
			So(len(groups), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("Create group", func() {
			message, err = conn.CreateGroup(group, owners, users)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("Group %s created.", group))
		})

		Convey("Add user to group", func() {
			message, err = conn.AddGroupUsers(group, owners, other)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("Group %s updated.", group))
		})

		Convey("Remove user from group", func() {
			message, err = conn.RemoveGroupUsers(group, owners, other)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("Group %s updated.", group))
		})

		Convey("Delete group", func() {
			message, err = conn.DeleteGroup(group)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("Group %s deleted.", group))
		})
	})
}
