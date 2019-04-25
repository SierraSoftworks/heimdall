package config

import (
	"testing"

	"github.com/SierraSoftworks/heimdall/internal/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Config", t, func() {
		confDir, err := test.GetTestAssetPath("config", "config_test")
		So(err, ShouldBeNil)
		So(confDir, ShouldNotEqual, "")

		t.Logf("Config Directory: %s", confDir)

		c, err := ReadConfig(confDir)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		So(c.Store, ShouldNotBeNil)
		So(c.Store.Type, ShouldEqual, "memory")
		So(c.Store.URL.String(), ShouldEqual, "memory://namespace")

		So(c.Publishers, ShouldHaveLength, 1)
		So(c.Publishers[0].Type, ShouldEqual, "redis")
		So(c.Publishers[0].URL.String(), ShouldEqual, "redis://localhost:6379")

		So(c.Collectors, ShouldHaveLength, 1)
		So(c.Collectors[0].Type, ShouldEqual, "nats")
		So(c.Collectors[0].URL.String(), ShouldEqual, "nats://localhost:4222")

		So(c.Checks, ShouldHaveLength, 2)
	})
}
