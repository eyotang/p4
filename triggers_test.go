package p4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggers_Triggers(t *testing.T) {
	var (
		currentTriggers []string
		lines           = []string{
			"swarm.commit change-commit //... \"%//.swarm/triggers/swarm-trigger.pl% -c %//.swarm/triggers/swarm-trigger.conf% -t commit -v %change%\"",
		}
	)
	conn, err := setup(t)
	Convey("test Triggers functions", t, func() {
		So(err, ShouldBeNil)

		Convey("List triggers", func() {
			currentTriggers, err = conn.Triggers()
			So(err, ShouldBeNil)
			So(len(currentTriggers), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("Preview triggers", func() {
			var (
				triggers = &Triggers{Lines: lines}
			)
			content := triggers.String()
			So(content, ShouldNotBeEmpty)
		})

		Convey("Write new triggers", func() {
			_, err = conn.WriteTriggers(lines)
			So(err, ShouldBeNil)
		})

		Convey("Restore old triggers", func() {
			_, err = conn.WriteTriggers(currentTriggers)
			So(err, ShouldBeNil)
		})
	})
}
