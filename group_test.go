package p4

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroup_Groups(t *testing.T) {
	var (
		message   string
		group     = "group-xxx"
		subGroups = []string{"eyotang2"}
		users     = []string{"eyotang", "tangyongqiang"}
		other     = []string{"abc"}
		owners    = []string{"owner"}
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
			message, err = conn.CreateGroup(group, owners, subGroups, users)
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

func TestGroup_GroupsRead(t *testing.T) {
	var (
		user  = "lejiajun"
		owner = "feihonghui"
	)
	conn, err := setup(t)
	Convey("test Group functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Display group", func() {
			var (
				group *GroupInfo
			)
			group, err = conn.GroupInfo("swarm-group")
			So(err, ShouldBeNil)
			So(len(group.Users), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("List my groups", func() {
			var (
				groups []string
			)

			groups, err = conn.GroupsBelong(user)
			So(err, ShouldBeNil)
			So(len(groups), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("List owned groups", func() {
			var (
				groups []string
			)

			groups, err = conn.GroupsOwned(owner)
			So(err, ShouldBeNil)
			So(len(groups), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
