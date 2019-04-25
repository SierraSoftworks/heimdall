package redis

import (
	"net/url"
	"testing"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/keimoon/gore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPublisher(t *testing.T) {
	Convey("RedisPublisher", t, func() {
		u, err := url.Parse("redis://localhost:6379/publisher")
		So(err, ShouldBeNil)

		d := &driver.Driver{
			Type: "redis",
			URL:  u,
		}

		p, err := NewRedisPublisher(d)
		So(err, ShouldBeNil)
		defer p.Close()

		Convey("Describe", func() {
			So(p.Describe(), ShouldEqual, "redis://localhost:6379/publisher")
		})

		Convey("Publish", func(c C) {
			rc, err := gore.Dial("localhost:6379")
			So(err, ShouldBeNil)
			defer rc.Close()

			subs := gore.NewSubscriptions(rc)
			So(err, ShouldBeNil)
			defer subs.Close()

			So(subs.Subscribe("/publisher/test"), ShouldBeNil)

			go func() {
				// Make sure we give enough time for the subscriber to start
				time.Sleep(50 * time.Millisecond)
				c.So(p.Publish("test", []byte("test")), ShouldBeNil)
			}()

			select {
			case m := <-subs.Message():
				So(m, ShouldNotBeNil)
				So(m.Channel, ShouldEqual, "/publisher/test")
				So(string(m.Message), ShouldEqual, "test")
			case <-time.After(time.Second):
				So("receive failed", ShouldBeNil)
			}
		})
	})
}
