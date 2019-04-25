package file

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCollector(t *testing.T) {
	Convey("FileCollector", t, func() {
		os.Remove("file_test.dat")
		defer os.Remove("file_test.dat")

		u, err := url.Parse("file://file_test.dat")
		So(err, ShouldBeNil)

		d := &driver.Driver{
			Type: "file",
			URL:  u,
		}

		c, err := NewFileCollector(d)
		So(err, ShouldBeNil)
		defer c.Close()

		Convey("Describe", func() {
			So(c.Describe(), ShouldEqual, "file://file_test.dat")
		})

		Convey("Subscribe", func() {
			s, err := c.Subscribe("test")
			So(err, ShouldBeNil)
			So(s, ShouldNotBeNil)

			defer s.Close()

			ch := s.Channel()
			So(ch, ShouldNotBeNil)

			fs, ok := s.(*fileSubscription)
			So(ok, ShouldBeTrue)

			f, err := os.OpenFile(fs.f.Name(), os.O_WRONLY|os.O_SYNC, 0664)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			f.WriteString("{\"channel\":\"test\",\"data\":\"dGVzdA==\"}\n")
			select {
			case msg := <-ch:
				So(string(msg), ShouldResemble, "test")
			case <-time.After(time.Second):
				So("receive failed", ShouldBeNil)
			}
		})
	})
}
