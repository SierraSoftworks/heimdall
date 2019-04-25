package file

import (
	"net/url"
	"os"
	"testing"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPublisher(t *testing.T) {
	Convey("FilePublisher", t, func() {
		os.Remove("file_test.dat")
		defer os.Remove("file_test.dat")

		u, err := url.Parse("file://file_test.dat")
		So(err, ShouldBeNil)

		d := &driver.Driver{
			Type: "file",
			URL:  u,
		}

		p, err := NewFilePublisher(d)
		So(err, ShouldBeNil)
		defer p.Close()

		Convey("Describe", func() {
			So(p.Describe(), ShouldEqual, "file://file_test.dat")
		})

		Convey("Publish", func() {
			fp, ok := p.(*FilePublisher)
			So(ok, ShouldBeTrue)

			f, err := os.OpenFile(fp.f.Name(), os.O_RDONLY|os.O_SYNC, 0664)
			So(err, ShouldBeNil)
			defer f.Close()

			So(fp.Publish("test", []byte("test")), ShouldBeNil)

			expected := "{\"channel\":\"test\",\"data\":\"dGVzdA==\"}\n"

			out := make([]byte, len(expected)*2)
			n, err := f.Read(out)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, len(expected))
			So(string(out[:n]), ShouldEqual, expected)
		})
	})
}
