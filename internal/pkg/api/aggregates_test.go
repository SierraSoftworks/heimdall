package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/duration"

	"net/http"
	"net/http/httptest"
	"net/url"

	"encoding/json"

	"github.com/SierraSoftworks/heimdall/internal/pkg/plugins/memory"
	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregates(t *testing.T) {
	Convey("Aggregates", t, func() {
		u, err := url.Parse("memory://memory")
		So(err, ShouldBeNil)

		store, err := plugins.GetStore(&driver.Driver{
			Type: "memory",
			URL:  u,
		})
		So(err, ShouldBeNil)
		So(store, ShouldNotBeNil)

		api, err := NewAPI(store)
		So(err, ShouldBeNil)

		s := httptest.NewServer(api.Router())
		defer s.Close()

		api.Store().(*memory.Memory).Reset()

		client := &models.Source{
			Name: "test-client",
			Tags: map[string]string{
				"type": "test",
			},
		}

		check := &models.Check{
			Name:    "apache-port80",
			Command: "curl -sS -D - http://localhost:80/",
			Collections: []string{
				"webservers",
			},
			Interval: duration.Duration(60 * time.Second),
			Timeout:  duration.Duration(5 * time.Second),
		}

		exec := &models.Execution{
			Scheduled: time.Now(),
			Executed:  time.Now(),
			Duration:  0,
			Status:    models.StatusOkay,
			Output:    "This is a quick test",
		}

		So(api.Store().AddReport(&models.Report{
			Source:    client,
			Check:     check,
			Execution: exec,
		}), ShouldBeNil)

		Convey("/api/v1/aggregates", func() {
			url := fmt.Sprintf("%s/api/v1/aggregates", s.URL)

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Source
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/aggregate/{aggregate}", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s", s.URL, "webservers")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data models.AggregateDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data.Name, ShouldEqual, "webservers")
			})

			Convey("DELETE", func() {
				r, err := http.NewRequest("DELETE", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data models.AggregateDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data.Name, ShouldEqual, "webservers")
			})
		})

		Convey("/api/v1/aggregate/{aggregate}/checks", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s/checks", s.URL, "webservers")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.CheckDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/aggregate/{aggregate}/clients", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s/clients", s.URL, "webservers")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Source
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})
	})
}
