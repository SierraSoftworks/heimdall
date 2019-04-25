package redis

import (
	"net/url"
	"testing"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/keimoon/gore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCollector(t *testing.T) {
	Convey("RedisCollector", t, func() {
		u, err := url.Parse("redis://localhost:6379/collector")
		So(err, ShouldBeNil)

		d := &driver.Driver{
			Type: "redis",
			URL:  u,
		}

		c, err := NewRedisCollector(d)
		So(err, ShouldBeNil)
		defer c.Close()

		Convey("Describe", func() {
			So(c.Describe(), ShouldEqual, "redis://localhost:6379/collector")
		})

		Convey("Subscribe", func() {
			s, err := c.Subscribe("test")
			So(err, ShouldBeNil)
			So(s, ShouldNotBeNil)

			defer s.Close()

			ch := s.Channel()
			So(ch, ShouldNotBeNil)

			rc, err := gore.Dial("localhost:6379")
			So(err, ShouldBeNil)
			defer rc.Close()

			go func() {
				// Give enough time for the subscriber to start
				time.Sleep(50 * time.Millisecond)
				gore.Publish(rc, "/collector/test", []byte("test"))
			}()

			select {
			case msg := <-ch:
				So(string(msg), ShouldEqual, "test")
			case <-time.After(time.Second):
				So("receive failed", ShouldBeNil)
			}
		})
	})
}
