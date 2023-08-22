package p4

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUser_Users(t *testing.T) {
	conn, err := setup(t)
	Convey("test Users functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List users", func() {
			users, err := conn.Users()
			So(err, ShouldBeNil)
			So(len(users), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("Get user info", func() {
			user, err := conn.User("readonly")
			So(err, ShouldBeNil)
			So(user.User, ShouldEqual, "readonly")
			So(user.AuthMethod, ShouldEqual, "perforce")
		})
	})
}

func TestUser_ListUsers(t *testing.T) {
	conn, err := setup(t)
	Convey("test Users functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List users without any group", func() {
			users, err := conn.Users()
			So(err, ShouldBeNil)
			So(len(users), ShouldBeGreaterThanOrEqualTo, 0)

			var unrefUsers []string
			for _, u := range users {
				groups, err := conn.GroupsBelong(u.User)
				So(err, ShouldBeNil)
				if len(groups) <= 0 {
					user, err := conn.User(u.User)
					So(err, ShouldBeNil)
					if user.AuthMethod == "ldap" {
						unrefUsers = append(unrefUsers, u.User)
					}
				}
			}
			So(len(unrefUsers), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestUser_DeleteUser(t *testing.T) {
	conn, err := setup(t)
	Convey("test User delete functions", t, func() {
		So(err, ShouldBeNil)

		Convey("Delete users without any group", func() {
			user := "localauth"
			message, err := conn.DeleteUser(user)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("Deletion of user %s and all the user's clients initiated.\nUser %s deleted.", user, user))
		})
	})
}

func TestUser_CreateUser(t *testing.T) {
	conn, err := setup(t)
	Convey("test User create functions", t, func() {
		So(err, ShouldBeNil)

		user := "test99"

		Convey("Delete user without any group", func() {
			_, err := conn.DeleteUser(user)
			So(err, ShouldBeNil)
		})

		Convey("Create user without any group", func() {
			userInfo := &UserInfo{
				User:       user,
				Email:      "tester99@mail.com",
				FullName:   "测试用户99",
				AuthMethod: "ldap",
			}
			message, err := conn.CreateUser(userInfo)
			So(err, ShouldBeNil)
			So(message, ShouldEqual, fmt.Sprintf("User %s saved.", user))
		})
	})
}
