package set

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSequenceSet(t *testing.T) {
	Convey("SequenceSet", t, func() {
		s := NewMemorySet()
		So(s, ShouldNotBeNil)

		has, err := s.Has("test")
		So(has, ShouldBeFalse)
		So(err, ShouldBeNil)

		So(s.Add("test", 1), ShouldBeNil)
		has, err = s.Has("test")
		So(has, ShouldBeTrue)
		So(err, ShouldBeNil)

		So(s.Add("test", 3), ShouldBeNil)
		has, err = s.Has("test")
		So(has, ShouldBeTrue)
		So(err, ShouldBeNil)

		So(s.Remove("test", 2), ShouldBeNil)
		has, err = s.Has("test")
		So(has, ShouldBeTrue)
		So(err, ShouldBeNil)

		So(s.Remove("test", 4), ShouldBeNil)
		has, err = s.Has("test")
		So(has, ShouldBeFalse)
		So(err, ShouldBeNil)

		So(s.Add("test", 5), ShouldBeNil)
		has, err = s.Has("test")
		So(has, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}
