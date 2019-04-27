package models

import (
	"strings"
	"testing"

	"bytes"
	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStatus(t *testing.T) {
	Convey("Status", t, func() {
		Convey("Constants", func() {
			So(StatusOkay, ShouldEqual, Status(0))
			So(StatusWarn, ShouldEqual, Status(1))
			So(StatusCrit, ShouldEqual, Status(2))
			So(StatusUnkn, ShouldEqual, Status(3))
		})

		Convey("String", func() {
			So(StatusOkay.String(), ShouldEqual, "OK")
			So(StatusWarn.String(), ShouldEqual, "WARN")
			So(StatusCrit.String(), ShouldEqual, "CRIT")
			So(StatusUnkn.String(), ShouldEqual, "UNKN")
			So(Status(127).String(), ShouldEqual, "UNKN")
		})

		Convey("ParseStatus", func() {
			So(ParseStatus("OK"), ShouldEqual, StatusOkay)
			So(ParseStatus("CRIT"), ShouldEqual, StatusCrit)
			So(ParseStatus("WARN"), ShouldEqual, StatusWarn)
			So(ParseStatus("UNKN"), ShouldEqual, StatusUnkn)
		})

		Convey("JSON", func() {
			Convey("Serialization", func() {
				cases := []struct {
					Status   Status `json:"status"`
					Expected string `json:"-"`
				}{
					{StatusOkay, `{"status":"OK"}`},
					{StatusWarn, `{"status":"WARN"}`},
					{StatusCrit, `{"status":"CRIT"}`},
					{StatusUnkn, `{"status":"UNKN"}`},
				}

				for _, c := range cases {
					b := bytes.NewBuffer([]byte{})
					So(json.NewEncoder(b).Encode(&c), ShouldBeNil)
					So(strings.TrimSpace(b.String()), ShouldEqual, c.Expected)
				}
			})

			Convey("Deserialization", func() {
				cases := []struct {
					Status  Status `json:"status"`
					Content string `json:"-"`
				}{
					{StatusOkay, `{"status":"OK"}`},
					{StatusWarn, `{"status":"WARN"}`},
					{StatusCrit, `{"status":"CRIT"}`},
					{StatusUnkn, `{"status":"UNKN"}`},
				}

				for _, c := range cases {
					var ct struct {
						Status Status `json:"status"`
					}
					b := bytes.NewBuffer([]byte(c.Content))
					So(json.NewDecoder(b).Decode(&ct), ShouldBeNil)
					So(ct.Status, ShouldEqual, c.Status)
				}
			})
		})
	})
}
