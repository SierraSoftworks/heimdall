package transports

import (
	"net/url"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisTransport(test *testing.T) {
	log.SetLevel(log.DebugLevel)
	defer log.SetLevel(log.ErrorLevel)

	Convey("RedisTransport", test, func() {
		u, err := url.Parse("redis://localhost:6379/heimdall")
		So(err, ShouldBeNil)
		So(u, ShouldNotBeNil)

		t, err := NewRedisTransport(u)
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		Convey("Describe", func() {
			So(t.Describe(), ShouldEqual, "redis://localhost:6379/heimdall")
		})

		Convey("Publish", func() {
			So(t.Publish("test/publish", []byte{}), ShouldBeNil)
		})

		Convey("Subscribe", func() {
			s, err := t.Subscribe("test/subscribe")
			So(err, ShouldBeNil)
			So(s, ShouldNotBeNil)

			time.Sleep(10 * time.Millisecond)

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
