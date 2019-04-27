package models

import (
	"testing"

	"bytes"
	"encoding/json"
	"strings"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDuration(t *testing.T) {
	Convey("Duration", t, func() {
		Convey("JSON", func() {
			type dS struct {
				D Duration `json:"d"`
			}

			Convey("Serialize", func() {
				d := dS{
					Duration(30 * time.Second),
				}
				b := bytes.NewBuffer([]byte{})

				So(json.NewEncoder(b).Encode(&d), ShouldBeNil)
				So(strings.TrimSpace(b.String()), ShouldEqual, `{"d":"30s"}`)
			})

			Convey("Deserialize", func() {
				var d dS
				b := bytes.NewBufferString(`{"d":"30s"}`)
				So(json.NewDecoder(b).Decode(&d), ShouldBeNil)
				So(d.D, ShouldEqual, Duration(30*time.Second))
			})
		})
	})
}
