package transports

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURL(t *testing.T) {
	Convey("URL", t, func() {
		Convey("nats://localhost:4222", func() {
			u, err := ParseURL("nats://localhost:4222")
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)

			So(u.Host, ShouldEqual, "localhost:4222")
			So(u.TopicPrefix, ShouldEqual, "")
			So(u.Options, ShouldResemble, map[string]string{})

			So(u.String(), ShouldEqual, "nats://localhost:4222")
			So(u.SafeString(), ShouldEqual, "nats://localhost:4222")
		})

		Convey("nats://localhost:4222/heimdall/custom", func() {
			u, err := ParseURL("nats://localhost:4222/heimdall/custom")
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)

			So(u.Host, ShouldEqual, "localhost:4222")
			So(u.TopicPrefix, ShouldEqual, "heimdall/custom")
			So(u.Options, ShouldResemble, map[string]string{})

			So(u.String(), ShouldEqual, "nats://localhost:4222/heimdall/custom")
			So(u.SafeString(), ShouldEqual, "nats://localhost:4222/heimdall/custom")

			So(u.GetFullTopic("test/publish"), ShouldEqual, "heimdall/custom/test/publish")
		})

		Convey("nats://localhost:4222?queue_group=prod", func() {
			u, err := ParseURL("nats://localhost:4222?queue_group=prod")
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)

			So(u.Host, ShouldEqual, "localhost:4222")
			So(u.TopicPrefix, ShouldEqual, "")
			So(u.Options, ShouldResemble, map[string]string{
				"queue_group": "prod",
			})

			So(u.String(), ShouldEqual, "nats://localhost:4222?queue_group=prod")
			So(u.SafeString(), ShouldEqual, "nats://localhost:4222")

			So(u.GetOption("queue_group", "heimdall_servers"), ShouldEqual, "prod")
			So(u.GetOption("timeout", "10"), ShouldEqual, "10")
		})

		Convey("nats://USER:PASS@localhost:4222", func() {
			u, err := ParseURL("nats://USER:PASS@localhost:4222")
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)

			So(u.Host, ShouldEqual, "localhost:4222")
			So(u.User, ShouldEqual, "USER")
			So(u.Password, ShouldEqual, "PASS")
			So(u.TopicPrefix, ShouldEqual, "")
			So(u.Options, ShouldResemble, map[string]string{})

			So(u.String(), ShouldEqual, "nats://USER:PASS@localhost:4222")
			So(u.SafeString(), ShouldEqual, "nats://localhost:4222")
		})

		Convey("nats://", func() {
			u, err := ParseURL("nats://")
			So(err, ShouldNotBeNil)
			So(u, ShouldBeNil)
		})

		Convey("file://./heimdall.dat", func() {
			u, err := ParseURL("file://./heimdall.dat")
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)

			So(u.Host, ShouldEqual, ".")
			So(u.TopicPrefix, ShouldEqual, "heimdall.dat")
			So(u.Options, ShouldResemble, map[string]string{})

			So(u.String(), ShouldEqual, "file://./heimdall.dat")
			So(u.SafeString(), ShouldEqual, "file://./heimdall.dat")
		})
	})
}
