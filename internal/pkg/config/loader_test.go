package config

import (
	"testing"

	"path/filepath"

	"github.com/SierraSoftworks/heimdall/internal/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigLoader(t *testing.T) {
	type Config struct {
		File    int  `json:"file"`
		Enabled bool `json:"enabled"`
		Extra   bool `json:"extra"`
	}

	Convey("Config Loader", t, func() {
		confDir, err := test.GetTestAssetPath("config", "loader_test")
		So(err, ShouldBeNil)
		So(confDir, ShouldNotEqual, "")

		Convey("FindConfig", func() {

			Convey("Directory", func() {
				configs, err := FindConfig(filepath.Join(confDir, "valid"))
				So(err, ShouldBeNil)
				So(configs, ShouldResemble, []string{
					filepath.Join(confDir, "valid", "config.yaml"),
					filepath.Join(confDir, "valid", "config.yml"),
				})

				Convey("Missing", func() {
					configs, err := FindConfig(filepath.Join(confDir, "nonexistent"))
					So(err, ShouldNotBeNil)
					So(configs, ShouldBeNil)
				})
			})

			Convey("File", func() {
				configs, err := FindConfig(filepath.Join(confDir, "valid", "config.yaml"))
				So(err, ShouldBeNil)
				So(configs, ShouldResemble, []string{
					filepath.Join(confDir, "valid", "config.yaml"),
				})

				Convey("Missing", func() {
					configs, err := FindConfig(filepath.Join(confDir, "nonexistent", "config.yaml"))
					So(err, ShouldNotBeNil)
					So(configs, ShouldBeNil)
				})
			})
		})

		Convey("LoadConfig", func() {
			c := Config{}

			So(LoadConfig(filepath.Join(confDir, "valid", "config.yaml"), &c), ShouldBeNil)
			So(c.File, ShouldEqual, 1)
			So(c.Enabled, ShouldBeTrue)
			So(c.Extra, ShouldBeFalse)

			So(LoadConfig(filepath.Join(confDir, "valid", "config.yml"), &c), ShouldBeNil)
			So(c.File, ShouldEqual, 2)
			So(c.Enabled, ShouldBeTrue)
			So(c.Extra, ShouldBeTrue)

			Convey("Directory", func() {
				err := LoadConfig(filepath.Join(confDir, "valid"), &c)
				So(err, ShouldNotBeNil)
			})

			Convey("Invalid", func() {
				err := LoadConfig(filepath.Join(confDir, "broken", "config.yaml"), &c)
				So(err, ShouldNotBeNil)
			})

			Convey("Missing", func() {
				err := LoadConfig(filepath.Join(confDir, "nonexistent"), &c)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
