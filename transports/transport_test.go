package transports

import (
	"testing"

	"net/url"
	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransports(t *testing.T) {
	Convey("Transports", t, func() {
		Convey("GetTransport", func() {
			Convey("NATS", func() {
				u, err := url.Parse("nats://localhost:4222")
				So(err, ShouldBeNil)

				tr, err := GetTransport("nats", u)
				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)
				defer tr.Close()

				So(tr, ShouldHaveSameTypeAs, &NATSTransport{})
			})

			Convey("Redis", func() {
				u, err := url.Parse("redis://localhost:6379")
				So(err, ShouldBeNil)

				tr, err := GetTransport("redis", u)
				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)
				defer tr.Close()

				So(tr, ShouldHaveSameTypeAs, &RedisTransport{})
			})

			Convey("File", func() {
				u, err := url.Parse("file://transport_test.dat")
				So(err, ShouldBeNil)

				tr, err := GetTransport("file", u)
				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)

				defer tr.Close()
				defer os.Remove("transport_test.dat")

				So(tr, ShouldHaveSameTypeAs, &FileTransport{})
			})
		})
	})
}
