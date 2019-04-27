package main

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Config", t, func() {
		cwd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(cwd, ShouldNotEqual, "")

		c, err := ReadConfig(filepath.Join(cwd, "examples"))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		So(c.Transports, ShouldHaveLength, 1)
		So(c.Transports[0].Driver, ShouldEqual, "nats")
		So(c.Transports[0].URL.String(), ShouldEqual, "nats://localhost:4222")

		So(c.Checks, ShouldHaveLength, 2)
	})
}
