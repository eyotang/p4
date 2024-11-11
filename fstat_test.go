package p4

import (
	"fmt"
	"os"
	"strings"
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

func TestConn_Fstat(t *testing.T) {
	var (
		conn *Conn
		err  error
	)

	conn, err = setup(t)
	Convey("test DescribeShelved functions", t, func() {
		So(err, ShouldBeNil)

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
		//So(message, ShouldEqual, fmt.Sprintf("Client %s saved.", client))
		So(err, ShouldBeNil)

		conn = conn.SetClient(client)
		var result []*Stat

		Convey("Describe Shelved", func() {
			result, err = conn.Fstats([]string{"//Arl.Private.Project/Mainline/main/Assets/77.txt"})
			So(result, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
