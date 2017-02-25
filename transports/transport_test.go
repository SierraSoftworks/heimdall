package transports

import (
	"testing"

	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransports(t *testing.T) {
	Convey("Transports", t, func() {
		Convey("GetTransport", func() {
			Convey("NATS", func() {
				tr, err := GetTransport("nats", "nats://localhost:4222")
				defer tr.Close()

				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)
				So(tr, ShouldHaveSameTypeAs, &NATSTransport{})
			})

			Convey("Redis", func() {
				tr, err := GetTransport("redis", "redis://localhost:6379")
				defer tr.Close()

				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)
				So(tr, ShouldHaveSameTypeAs, &RedisTransport{})
			})

			Convey("File", func() {
				tr, err := GetTransport("file", "file://transport_test.dat")
				defer tr.Close()
				defer os.Remove("transport_test.dat")

				So(err, ShouldBeNil)
				So(tr, ShouldNotBeNil)
				So(tr, ShouldHaveSameTypeAs, &FileTransport{})
			})
		})
	})
}
