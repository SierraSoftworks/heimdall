package driver

import (
	"net/url"
	"testing"

	"bytes"
	"encoding/json"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDriver(t *testing.T) {
	Convey("Driver", t, func() {
		Convey("JSON", func() {
			u, err := url.Parse("test://user:pass@host/path?opt=test")
			So(err, ShouldBeNil)

			driver := Driver{
				Type: "test",
				URL:  u,
			}

			Convey("Describe", func() {
				Convey("With just a type", func() {
					So((&Driver{Type: "test"}).Describe(), ShouldEqual, "test")
				})

				Convey("With a type and URL", func() {
					So((&Driver{Type: "test", URL: u}).Describe(), ShouldEqual, "test(test://user:pass@host/path?opt=test)")
				})

				Convey("With a type, URL and filter", func() {
					So((&Driver{Type: "test", URL: u, Filter: map[string]interface{}{"x": 1}}).Describe(), ShouldEqual, `test(test://user:pass@host/path?opt=test) {"x":1}`)
				})
			})

			Convey("Matches", func() {
				Convey("With no filter", func() {
					d := Driver{}
					So(d.Matches(map[string]interface{}{}), ShouldBeTrue)
					So(d.Matches(map[string]interface{}{
						"test": true,
					}), ShouldBeTrue)
				})

				Convey("With a filter", func() {
					d := Driver{
						Filter: map[string]interface{}{
							"test": true,
						},
					}

					So(d.Matches(map[string]interface{}{
						"test": true,
					}), ShouldBeTrue)

					So(d.Matches(map[string]interface{}{
						"test": false,
					}), ShouldBeFalse)
				})
			})

			Convey("Equals", func() {
				So(driver.Equals(&Driver{
					Type: "test",
					URL:  u,
				}), ShouldBeTrue)
			})

			Convey("Serialize", func() {
				b := bytes.NewBuffer([]byte{})

				So(json.NewEncoder(b).Encode(&driver), ShouldBeNil)
				So(strings.TrimSpace(b.String()), ShouldEqual, `{"type":"test","url":"test://user:pass@host/path?opt=test"}`)
			})

			Convey("Deserialize", func() {
				var d Driver
				b := bytes.NewBufferString(`{"type":"test", "url":"test://user:pass@host","filter":{"x":1}}`)
				So(json.NewDecoder(b).Decode(&d), ShouldBeNil)
				So(d.Type, ShouldEqual, "test")
				So(d.URL.String(), ShouldEqual, "test://user:pass@host")
				So(d.Filter, ShouldHaveSameTypeAs, map[string]interface{}{})
				So(d.Filter["x"], ShouldEqual, 1)

				b = bytes.NewBufferString(`{"type":"test", "url":"test://user:pass@host","filter":{0:1}}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)

				b = bytes.NewBufferString(`{"type":"test", "url":"test://user:pass@host","filter":[0]}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)

				b = bytes.NewBufferString(`{"type":"test", "url":"test://user:pass@host","filter":"test"}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)
			})

			Convey("SafeURLString", func() {
				So(driver.SafeURLString(), ShouldEqual, "test://host/path")
			})

			Convey("SafeURLHost", func() {
				So(driver.SafeURLHost(), ShouldEqual, "test://host")
			})

			Convey("GetPath", func() {
				So(driver.GetPath("test/channel"), ShouldEqual, "/path/test/channel")
				So(driver.GetPath("/test/channel"), ShouldEqual, "/path/test/channel")
			})

			Convey("GetOption", func() {
				So(driver.GetOption("opt", "not found"), ShouldEqual, "test")
				So(driver.GetOption("opt2", "default"), ShouldEqual, "default")
			})
		})
	})
}
