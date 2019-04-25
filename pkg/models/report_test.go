package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReport(t *testing.T) {
	Convey("Report", t, func() {
		r := &Report{
			UUID: "test",
			Source: &Source{
				UUID: "test",
				Name: "Test Client",
			},
			Check: &Check{
				UUID:    "test",
				Command: "/bin/false",
			},
			Execution: &Execution{
				Status: StatusOkay,
			},
		}

		Convey("ToMap()", func() {
			m := r.ToMap()
			So(m, ShouldNotBeNil)
			So(m, ShouldHaveSameTypeAs, map[string]interface{}{})
			So(m, ShouldContainKey, "uuid")
			So(m, ShouldContainKey, "source")
			So(m, ShouldContainKey, "check")
			So(m, ShouldContainKey, "execution")
		})
	})
}
