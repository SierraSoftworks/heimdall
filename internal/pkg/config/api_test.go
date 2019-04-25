package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {
	Convey("API", t, func() {
		a := &APIConfig{
			Listen: ":8080",
		}

		So(a, ShouldNotBeNil)
		So(a.Listen, ShouldEqual, ":8080")

		a.Update(&APIConfig{
			Listen: "127.0.0.1:8080",
		})
		So(a.Listen, ShouldEqual, "127.0.0.1:8080")
	})
}
