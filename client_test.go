package p4

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient_Clients(t *testing.T) {
	var (
		conn *Conn
		err  error
	)
	conn, err = setup(t)
	Convey("test Clients", t, func() {
		So(err, ShouldBeNil)

		Convey("List clients", func() {
			clients, err := conn.Clients("//DM99.ZGame.Project/Main/ZGame_Mainline")
			So(clients, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
	})
}

func TestClient_Client(t *testing.T) {
	var (
		err  error
		conn *Conn
	)

	conn, err = setup(t)
	Convey("test Client", t, func() {
		So(err, ShouldBeNil)

		Convey("Display client", func() {
			client, err := conn.Client("root_ZGame_Mainline")
			So(client, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Create partitioned client and delete", func() {
			owner := os.Getenv("P4USER")
			client := owner + "_xxx_ZGame_Mainline"
			wsRoot, _ := os.Getwd()
			clientInfo := Client{
				Client:        client,
				Owner:         owner,
				Host:          "",
				Description:   "",
				Root:          wsRoot + "/" + client,
				Options:       "noallwrite noclobber nocompress unlocked nomodtime normdir",
				SubmitOptions: "submitunchanged",
				Stream:        "",
				Type:          "",
				View:          []string{fmt.Sprintf("//DM99.ZGame.Project/Main/ZGame_Mainline/... //%s/...", client)},
			}
			message, err := conn.CreatePartitionClient(clientInfo)
			So(message, ShouldNotBeEmpty)
			So(err, ShouldBeNil)

			message, err = conn.DeleteClient(client)
			So(message, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
	})
}
