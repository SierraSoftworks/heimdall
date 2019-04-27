package utils

import (
	"testing"

	"os"
	"path/filepath"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	type Config struct {
		File    int  `json:"file"`
		Enabled bool `json:"enabled"`
		Extra   bool `json:"extra"`
	}

	Convey("Config", t, func() {
		cwd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(cwd, ShouldNotEqual, "")

		Convey("FindConfig", func() {

			Convey("Directory", func() {
				configs, err := FindConfig(filepath.Join(cwd, "examples", "valid"))
				So(err, ShouldBeNil)
				So(configs, ShouldResemble, []string{
					filepath.Join(cwd, "examples", "valid", "config.yaml"),
					filepath.Join(cwd, "examples", "valid", "config.yml"),
				})

				Convey("Missing", func() {
					configs, err := FindConfig(filepath.Join(cwd, "examples", "nonexistent"))
					So(err, ShouldNotBeNil)
					So(configs, ShouldBeNil)
				})
			})

			Convey("File", func() {
				configs, err := FindConfig(filepath.Join(cwd, "examples", "valid", "config.yaml"))
				So(err, ShouldBeNil)
				So(configs, ShouldResemble, []string{
					filepath.Join(cwd, "examples", "valid", "config.yaml"),
				})

				Convey("Missing", func() {
					configs, err := FindConfig(filepath.Join(cwd, "examples", "nonexistent", "config.yaml"))
					So(err, ShouldNotBeNil)
					So(configs, ShouldBeNil)
				})
			})
		})

		Convey("LoadConfig", func() {
			c := Config{}

			So(LoadConfig(filepath.Join(cwd, "examples", "valid", "config.yaml"), &c), ShouldBeNil)
			So(c.File, ShouldEqual, 1)
			So(c.Enabled, ShouldBeTrue)
			So(c.Extra, ShouldBeFalse)

			So(LoadConfig(filepath.Join(cwd, "examples", "valid", "config.yml"), &c), ShouldBeNil)
			So(c.File, ShouldEqual, 2)
			So(c.Enabled, ShouldBeTrue)
			So(c.Extra, ShouldBeTrue)

			Convey("Directory", func() {
				err := LoadConfig(filepath.Join(cwd, "examples", "valid"), &c)
				So(err, ShouldNotBeNil)
			})

			Convey("Invalid", func() {
				err := LoadConfig(filepath.Join(cwd, "examples", "broken", "config.yaml"), &c)
				So(err, ShouldNotBeNil)
			})

			Convey("Missing", func() {
				err := LoadConfig(filepath.Join(cwd, "examples", "nonexistent"), &c)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
