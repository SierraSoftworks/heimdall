package config

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeepalive(t *testing.T) {
	Convey("Keepalive", t, func() {
		k := &KeepaliveConfig{
			Aggregates: []string{},
			Interval:   30 * time.Second,
		}

		So(k, ShouldNotBeNil)
		So(k.Aggregates, ShouldResemble, []string{})
		So(k.Interval, ShouldEqual, 30*time.Second)

		k.Update(&KeepaliveConfig{
			Interval: 60 * time.Second,
		})

		So(k.Interval, ShouldEqual, 60*time.Second)
		So(k.Aggregates, ShouldResemble, []string{})

		k.Update(&KeepaliveConfig{
			Aggregates: []string{"test"},
		})

		So(k.Interval, ShouldEqual, 60*time.Second)
		So(k.Aggregates, ShouldResemble, []string{"test"})
	})
}
