package transports

import (
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileTransport(test *testing.T) {

	Convey("FileTransport", test, func() {
		os.Remove("file_test.dat")

		u, err := ParseURL("file://file_test.dat")
		So(err, ShouldBeNil)
		So(u, ShouldNotBeNil)

		t, err := NewFileTransport(u)
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		Convey("Describe", func() {
			So(t.Describe(), ShouldEqual, "file://file_test.dat")
		})

		Convey("Publish", func() {
			So(t.Publish("test/publish", []byte{}), ShouldBeNil)
		})

		Convey("Subscribe", func() {
			log.SetLevel(log.DebugLevel)
			defer log.SetLevel(log.ErrorLevel)

			s, err := t.Subscribe("test/subscribe")
			So(err, ShouldBeNil)
			So(s, ShouldNotBeNil)
			defer s.Close()

			time.Sleep(time.Second)

			t.Publish("test/subscribe", []byte("testing"))

			log.Debug("Reading from channel")
			select {
			case data := <-s.Channel():
				So(string(data), ShouldEqual, "testing")
			case <-time.After(time.Second):
				So("receive failed", ShouldBeNil)
			}

			time.Sleep(time.Second)

			So(s.Close(), ShouldBeNil)
		})

		So(t.Close(), ShouldBeNil)
		So(os.Remove("file_test.dat"), ShouldBeNil)

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
