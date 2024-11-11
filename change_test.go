package p4

import (
	"fmt"
	"os"
	"strings"
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

func TestChange_NewChangeList(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test changes", t, func() {
		So(err, ShouldBeNil)

		Convey("List shelved", func() {
			// 查询stream
			stream, err := conn.ChangeListStream(20585)
			So(stream, ShouldNotBeEmpty)
			So(err, ShouldBeNil)

			// 创建临时partitioned workspace
			streamWs := strings.Trim(stream, "/")

			// root_DM99.ZGame.Project-Development-xiner_test
			client := "root" + "_" + strings.ReplaceAll(streamWs, "/", "-")
			wsRoot, _ := os.Getwd()
			clientInfo := Client{
				Client:        client,
				Owner:         "root",
				Root:          wsRoot + "/" + client,
				Options:       "noallwrite noclobber nocompress unlocked nomodtime normdir",
				SubmitOptions: "submitunchanged",
				Stream:        stream,
				View:          []string{fmt.Sprintf("%s/... //%s/...", stream, client)},
			}
			message, err := conn.CreatePartitionClient(clientInfo)
			So(message, ShouldNotBeEmpty)
			// Client root_DM99.ZGame.Project-Development-xiner_test saved.
			So(message, ShouldEqual, fmt.Sprintf("Client %s saved.", client))
			So(err, ShouldBeNil)

			conn = conn.SetClient(client)
			change, err := conn.NewChangeList(NewChangeList{
				Change:      "new",
				User:        "sunqi01",
				Description: "123",
			})
			So(change, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestConn_DeleteChange(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test changes", t, func() {
		So(err, ShouldBeNil)

		Convey("List shelved", func() {
			conn = conn.SetClient("root_Arl.Private.Project-Mainline-main1")

			conn.ChangeUser("sunqi01", "B2697CD7CC377C6AB86CA886B09E81CA")
			message, err := conn.DeleteChange(17529)
			fmt.Println(message)
			So(message, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
