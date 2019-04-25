package test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAssets(t *testing.T) {
	Convey("Assets", t, func() {
		Convey("GetRepoRoot", func() {
			repo, err := GetRepoRoot()
			So(err, ShouldBeNil)
			So(repo, ShouldNotEqual, "")

			stat, err := os.Stat(repo)
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			So(stat.IsDir(), ShouldBeTrue)
		})

		Convey("GetAssetsPath", func() {
			assets, err := GetAssetsPath()
			So(err, ShouldBeNil)
			So(assets, ShouldNotEqual, "")
			So(assets, ShouldEndWith, "assets")

			stat, err := os.Stat(filepath.Join(assets, ".gitkeep"))
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			So(stat.IsDir(), ShouldBeFalse)
		})

		Convey("GetAssetPath", func() {
			asset, err := GetAssetPath(".gitkeep")
			So(err, ShouldBeNil)
			So(asset, ShouldNotEqual, "")
			So(asset, ShouldEndWith, filepath.Join("assets", ".gitkeep"))

			stat, err := os.Stat(asset)
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			So(stat.IsDir(), ShouldBeFalse)
		})

		Convey("GetTestAssetsPath", func() {
			assets, err := GetTestAssetsPath()
			So(err, ShouldBeNil)
			So(assets, ShouldNotEqual, "")
			So(assets, ShouldEndWith, "test")

			stat, err := os.Stat(filepath.Join(assets, ".gitkeep"))
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			So(stat.IsDir(), ShouldBeFalse)
		})

		Convey("GetTestAssetPath", func() {
			asset, err := GetTestAssetPath(".gitkeep")
			So(err, ShouldBeNil)
			So(asset, ShouldNotEqual, "")
			So(asset, ShouldEndWith, filepath.Join("test", ".gitkeep"))

			stat, err := os.Stat(asset)
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			So(stat.IsDir(), ShouldBeFalse)
		})
	})
}
