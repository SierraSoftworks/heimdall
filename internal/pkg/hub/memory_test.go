package hub

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryHub(t *testing.T) {
	Convey("MemoryHub", t, func(c C) {
		hub := NewMemoryHub()
		So(hub, ShouldNotBeNil)

		memHub, ok := hub.(*memoryHub)
		So(ok, ShouldBeTrue)

		msg := 1

		handled := make(chan struct{})
		sub := NewCallbackSubscriber(func(m interface{}) {
			handled <- struct{}{}
			c.So(m, ShouldNotBeNil)
			c.So(m, ShouldEqual, msg)
		})

		hub.Subscribe(sub)

		Convey("Subscribe", func() {
			So(memHub.subscribers, ShouldHaveLength, 1)

			Convey("Duplicates", func() {
				hub.Subscribe(sub)
				So(memHub.subscribers, ShouldHaveLength, 1)
			})
		})

		Convey("Notify", func() {
			hub.Notify(msg)
			select {
			case <-handled:
			case <-time.After(time.Second):
				So("not handled", ShouldBeNil)
			}
		})

		Convey("Unsubscribe", func() {
			hub.Unsubscribe(sub)
			So(memHub.subscribers, ShouldHaveLength, 0)
		})
	})
}
