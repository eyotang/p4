package p4

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubmit_SubmitShelve(t *testing.T) {
	var (
		err  error
		conn *Conn
	)

	conn, err = setup(t)
	Convey("test Submit", t, func() {
		So(err, ShouldBeNil)

		Convey("Create partitioned client, reshelve, submit and delete client", func() {
			owner := os.Getenv("P4USER")
			var (
				shelve     uint64 = 8113
				reshelveCL uint64
				submitCL   uint64
				submitter  = "tangyongqiang"
			)

			// 查询stream
			stream, err := conn.ChangeListStream(shelve)
			So(stream, ShouldNotBeEmpty)
			So(err, ShouldBeNil)

			// 创建临时partitioned workspace
			streamWs := strings.Trim(stream, "/")

			// root_DM99.ZGame.Project-Development-xiner_test
			client := owner + "_" + strings.ReplaceAll(streamWs, "/", "-")
			wsRoot, _ := os.Getwd()
			clientInfo := Client{
				Client:        client,
				Owner:         owner,
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

			// 将shelve CL给reshelve成新的shelve CL
			conn = conn.SetClient(client)
			message, err = conn.Reshelve(shelve)
			So(message, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
			regex := regexp.MustCompile("Change (\\d+) files shelved.")
			fields := regex.FindStringSubmatch(message)
			if len(fields) == 2 {
				reshelveCL, _ = strconv.ParseUint(fields[1], 10, 64)
			}

			// 提交新的shelve CL
			message, err = conn.SubmitShelve(reshelveCL)
			So(message, ShouldNotBeEmpty)
			if err != nil {
				message, err = conn.DeleteShelve(reshelveCL)
				So(message, ShouldNotBeEmpty)
				So(err, ShouldBeNil)
			} else {
				// 提交成功，修改提交人
				regex1 := regexp.MustCompile("renamed change (\\d+) and submitted.")
				regex2 := regexp.MustCompile("Change (\\d+) submitted.")
				fields = regex1.FindStringSubmatch(message)
				if len(fields) == 2 {
					submitCL, _ = strconv.ParseUint(fields[1], 10, 64)
				} else {
					fields = regex2.FindStringSubmatch(message)
					if len(fields) == 2 {
						submitCL, _ = strconv.ParseUint(fields[1], 10, 64)
					}
				}
				cl, err := conn.ChangeList(submitCL)
				So(err, ShouldBeNil)

				cl.User = submitter
				message, err = conn.UpdateChangeList(*cl)
				So(message, ShouldEqual, fmt.Sprintf("Change %d updated.", submitCL))
				So(err, ShouldBeNil)
			}

			// 删除临时partitioned workspace
			message, err = conn.DeleteClient(client)
			So(message, ShouldNotBeEmpty)
			// Client root_DM99.ZGame.Project-Development-xiner_test deleted.
			So(message, ShouldEqual, fmt.Sprintf("Client %s deleted.", client))
			So(err, ShouldBeNil)
		})
	})
}
