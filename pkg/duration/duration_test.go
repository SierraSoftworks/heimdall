package duration

import (
	"testing"

	"bytes"
	"encoding/json"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

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

				b = bytes.NewBufferString(`{"d":"30"}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)

				b = bytes.NewBufferString(`{"d":30}`)
				So(json.NewDecoder(b).Decode(&d), ShouldBeNil)
				So(d.D, ShouldEqual, Duration(30*time.Second))

				b = bytes.NewBufferString(`{"d":3}`)
				So(json.NewDecoder(b).Decode(&d), ShouldBeNil)
				So(d.D, ShouldEqual, Duration(3*time.Second))

				b = bytes.NewBufferString(`{"d":thirty}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)

				b = bytes.NewBufferString(`{"d":true}`)
				So(json.NewDecoder(b).Decode(&d), ShouldNotBeNil)
			})
		})

		Convey("YAML", func() {
			type dS struct {
				D Duration `json:"d"`
			}

			Convey("Serialize", func() {
				b, err := yaml.Marshal(&dS{Duration(30 * time.Second)})
				So(err, ShouldBeNil)
				So(strings.TrimSpace(string(b)), ShouldEqual, `d: 30s`)
			})

			Convey("Deserialize", func() {
				var d dS
				So(yaml.Unmarshal([]byte("d: 30s"), &d), ShouldBeNil)
				So(d.D, ShouldEqual, Duration(30*time.Second))
			})
		})
	})
}
