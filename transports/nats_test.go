package transports

import (
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNATSTransport(test *testing.T) {
	Convey("NATSTransport", test, func() {
		u, err := url.Parse("nats://localhost:4222/heimdall")
		So(err, ShouldBeNil)
		So(u, ShouldNotBeNil)

		t, err := NewNATSTransport(u)
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		Convey("Describe", func() {
			So(t.Describe(), ShouldEqual, "nats://localhost:4222/heimdall")
		})

		Convey("Publish", func() {
			So(t.Publish("test/publish", []byte{}), ShouldBeNil)
		})

		Convey("Subscribe", func() {
			s, err := t.Subscribe("test/subscribe")
			So(err, ShouldBeNil)
			So(s, ShouldNotBeNil)

			t.Publish("test/subscribe", []byte("testing"))

			select {
			case data := <-s.Channel():
				So(string(data), ShouldEqual, "testing")
			case <-time.After(time.Second):
				So("receive failed", ShouldBeNil)
			}

			So(s.Close(), ShouldBeNil)
		})

		So(t.Close(), ShouldBeNil)

		Convey("After Closing", func() {
			Convey("Publish", func() {
				So(t.Publish("test/publish", []byte{}), ShouldNotBeNil)
			})

			Convey("Subscribe", func() {
				s, err := t.Subscribe("test/subscribe")
				So(err, ShouldNotBeNil)
				So(s, ShouldBeNil)
			})
		})
	})
}
