package lamport

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClock(t *testing.T) {
	Convey("Lamport Clock", t, func() {
		cl := NewClock(0)

		// Next should increment the current time
		So(cl.Next(), ShouldEqual, uint64(1))
		So(cl.Next(), ShouldEqual, uint64(2))

		// Update with a higher timestamp should increment past that
		So(cl.Update(10), ShouldEqual, uint64(11))

		// Update with a lower timestamp should increment the current time
		So(cl.Update(5), ShouldEqual, uint64(12))

		// Current shouldn't modify the timestamp
		So(cl.Current(), ShouldEqual, uint64(12))

		// NotAfter should work correctly for different boundaries
		So(cl.NotAfter(1), ShouldBeFalse)
		So(cl.NotAfter(12), ShouldBeFalse)
		So(cl.NotAfter(15), ShouldBeTrue)
	})
}
