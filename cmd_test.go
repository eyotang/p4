package p4

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConn_Timeout(t *testing.T) {
	var (
		conn *Conn
		err  error
	)

	conn, err = setup(t)
	Convey("test Timeout functions", t, func() {
		So(err, ShouldBeNil)

		// 登录另外错误端口
		address := os.Getenv("P4PORT2")
		timeout := &Timeout{Login: 5 * time.Second}
		conn, err = NewConn(address, conn.username, conn.password, WithTimeout(timeout))
		So(err, ShouldNotBeNil)

		address = os.Getenv("P4PORT2")
		conn, err = NewConn(address, conn.username, conn.password)
		So(err, ShouldNotBeNil)
	})
}
